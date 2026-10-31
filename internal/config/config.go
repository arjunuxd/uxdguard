package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config represents the project configuration.
type Config struct {
	Version  int              `yaml:"version"`
	Scan     ScanConfig       `yaml:"scan"`
	Policy   PolicyConfig     `yaml:"policy"`
	Rules    RulesConfig      `yaml:"rules,omitempty"`
	External ExternalScanners `yaml:"external,omitempty"`
}

// ExternalScanners configures external tools.
type ExternalScanners struct {
	Gitleaks ExternalConfig `yaml:"gitleaks,omitempty"`
	Semgrep  ExternalConfig `yaml:"semgrep,omitempty"`
}

// RulesConfig configures rules.
type RulesConfig struct {
	Enable   []string          `yaml:"enable,omitempty"`
	Disable  []string          `yaml:"disable,omitempty"`
	Severity map[string]string `yaml:"severity,omitempty"`
}

// ScanConfig contains scan settings.
type ScanConfig struct {
	MaxFileSizeBytes int64    `yaml:"max_file_size_bytes"`
	FollowSymlinks   bool     `yaml:"follow_symlinks"`
	Exclusions       []string `yaml:"exclusions"`
}

// PolicyConfig contains policy settings.
type PolicyConfig struct {
	FailOn []string `yaml:"fail_on"`
}

// Load loads configuration from the given path. If path is empty,
// it looks for .uxdguard.yaml in the project root. If not found,
// returns defaults.
func Load(projectRoot, configPath string) (*Config, string, error) {
	if projectRoot == "" {
		var err error
		projectRoot, err = os.Getwd()
		if err != nil {
			return nil, "", err
		}
	}
	absRoot, err := filepath.Abs(projectRoot)
	if err != nil {
		return nil, "", err
	}
	path := configPath
	if path == "" {
		path = filepath.Join(absRoot, ".uxdguard.yaml")
	} else {
		path, err = filepath.Abs(path)
		if err != nil {
			return nil, "", err
		}
		// Prevent a config file outside the project root, including paths on a
		// different volume where a relative path cannot be computed.
		rel, err := filepath.Rel(absRoot, path)
		if err != nil || rel == ".." || strings.HasPrefix(filepath.ToSlash(rel), "../") {
			return nil, path, errors.New("config file must be within project root")
		}
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if configPath != "" {
			return nil, path, fmt.Errorf("config file not found: %s", path)
		}
		// An absent project config is fine; use defaults.
		cfg := Defaults()
		return cfg, path, nil
	} else if err != nil {
		return nil, path, err
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, path, err
	}
	if len(b) == 0 {
		return nil, path, errors.New("empty configuration file")
	}
	var raw map[string]interface{}
	if err := yaml.Unmarshal(b, &raw); err != nil {
		return nil, path, fmt.Errorf("malformed YAML in %s: %w", path, err)
	}
	// Also try to unmarshal to struct
	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, path, fmt.Errorf("malformed YAML in %s: %w", path, err)
	}
	// Merge defaults only for keys that are absent. Explicit zero or empty
	// values must reach validation instead of being silently replaced.
	def := Defaults()
	if !hasKey(raw, "version") {
		cfg.Version = def.Version
	}
	if !hasNestedKey(raw, "scan", "max_file_size_bytes") {
		cfg.Scan.MaxFileSizeBytes = def.Scan.MaxFileSizeBytes
	}
	if !hasNestedKey(raw, "scan", "exclusions") {
		cfg.Scan.Exclusions = append([]string(nil), def.Scan.Exclusions...)
	}
	if !hasNestedKey(raw, "policy", "fail_on") {
		cfg.Policy.FailOn = append([]string(nil), def.Policy.FailOn...)
	}
	if err := Validate(&cfg); err != nil {
		return nil, path, err
	}
	return &cfg, path, nil
}

func hasKey(m map[string]interface{}, key string) bool {
	_, ok := m[key]
	return ok
}

func hasNestedKey(m map[string]interface{}, parent, key string) bool {
	section, ok := m[parent].(map[string]interface{})
	if !ok {
		return false
	}
	_, ok = section[key]
	return ok
}

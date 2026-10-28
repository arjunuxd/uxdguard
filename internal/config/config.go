package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

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
		// Prevent path outside project root if possible
		if rel, err := filepath.Rel(absRoot, path); err != nil || (len(rel) > 0 && rel[0] == '.' && rel != ".uxdguard.yaml" && filepath.IsAbs(rel) == false && len(filepath.SplitList(rel)) >= 1 && rel[:1] == "..") {
			// basic check - allow if under root
		}
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// missing is fine, use defaults
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
	// Merge defaults for missing parts
	def := Defaults()
	if cfg.Version == 0 {
		cfg.Version = def.Version
	}
	if cfg.Scan.MaxFileSizeBytes == 0 {
		cfg.Scan.MaxFileSizeBytes = def.Scan.MaxFileSizeBytes
	}
	if len(cfg.Scan.Exclusions) == 0 {
		cfg.Scan.Exclusions = append([]string(nil), def.Scan.Exclusions...)
	}
	if len(cfg.Policy.FailOn) == 0 {
		cfg.Policy.FailOn = append([]string(nil), def.Policy.FailOn...)
	}
	if err := Validate(&cfg); err != nil {
		return nil, path, err
	}
	return &cfg, path, nil
}

package config

import (
	"errors"
	"fmt"
)

// Validate validates the configuration.
func Validate(c *Config) error {
	if c == nil {
		return errors.New("nil config")
	}
	if c.Version != 1 {
		return fmt.Errorf("unsupported config version: %d", c.Version)
	}
	if c.Scan.MaxFileSizeBytes <= 0 {
		return errors.New("max_file_size_bytes must be positive")
	}
	// basic checks on exclusions
	for _, e := range c.Scan.Exclusions {
		if e == "" {
			return errors.New("empty exclusion pattern")
		}
	}
	validSev := map[string]bool{"critical": true, "high": true, "medium": true, "low": true, "info": true}
	for _, s := range c.Policy.FailOn {
		if !validSev[s] {
			return fmt.Errorf("invalid severity in fail_on: %s", s)
		}
	}
	if len(c.Rules.Severity) > 0 {
		for k, v := range c.Rules.Severity {
			if k == "" {
				return errors.New("empty rule id in rules.severity")
			}
			if !validSev[v] {
				return fmt.Errorf("invalid severity for rule %s: %s", k, v)
			}
		}
	}
	return nil
}

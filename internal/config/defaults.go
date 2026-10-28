package config

// Defaults returns the default configuration.
func Defaults() *Config {
	return &Config{
		Version: 1,
		Scan: ScanConfig{
			MaxFileSizeBytes: 1024 * 1024,
			FollowSymlinks:   false,
			Exclusions: []string{
				".git/**",
				"vendor/**",
				"node_modules/**",
				"dist/**",
				"build/**",
			},
		},
		Policy: PolicyConfig{
			FailOn: []string{"critical", "high"},
		},
	}
}

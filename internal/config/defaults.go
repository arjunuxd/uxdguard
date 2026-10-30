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
				".next/**",
				".nuxt/**",
				".svelte-kit/**",
				".turbo/**",
				".cache/**",
				"coverage/**",
				"target/**",
				"out/**",
				".venv/**",
				"venv/**",
				"__pycache__/**",
				".terraform/**",
				".idea/**",
				".vscode/**",
			},
		},
		Policy: PolicyConfig{
			FailOn: []string{"critical", "high"},
		},
	}
}

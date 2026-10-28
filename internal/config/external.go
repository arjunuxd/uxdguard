package config

// ExternalConfig configures external scanners.
type ExternalConfig struct {
	Enabled bool              `yaml:"enabled,omitempty"`
	Require bool              `yaml:"require,omitempty"`
	Path    string            `yaml:"path,omitempty"`
	Timeout string            `yaml:"timeout,omitempty"`
	Options map[string]string `yaml:"options,omitempty"`
}

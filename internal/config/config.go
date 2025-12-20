package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Agent       AgentConfig       `yaml:"agent"`
	Exporter    ExporterConfig    `yaml:"exporter"`
	Modules     ModulesConfig     `yaml:"modules"`
	Remediation RemediationConfig `yaml:"remediation"`
}

type AgentConfig struct {
	Name      string `yaml:"name"`
	LogLevel  string `yaml:"log_level"`
	ServerURL string `yaml:"server_url"`
}

type ExporterConfig struct {
	Splunk SplunkConfig `yaml:"splunk"`
}

type SplunkConfig struct {
	Enabled       bool          `yaml:"enabled"`
	URL           string        `yaml:"url"`
	Token         string        `yaml:"token"`
	BatchInterval time.Duration `yaml:"batch_interval"`
}

type ModulesConfig struct {
	Osquery OsqueryConfig `yaml:"osquery"`
}

type OsqueryConfig struct {
	Enabled    bool   `yaml:"enabled"`
	BinaryPath string `yaml:"binary_path"`
}

type RemediationConfig struct {
	Enabled           bool     `yaml:"enabled"`
	AutoPatchSeverity []string `yaml:"auto_patch_severity"`
}

// LoadConfig reads and parses the YAML configuration file.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

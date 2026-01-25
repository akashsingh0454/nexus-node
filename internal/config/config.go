package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Agent       AgentConfig       `yaml:"agent"`
	Security    SecurityConfig    `yaml:"security"`
	Pipeline    PipelineConfig    `yaml:"pipeline"`
	Modules     ModulesConfig     `yaml:"modules"`
	Remediation RemediationConfig `yaml:"remediation"`
}

type AgentConfig struct {
	Name      string `yaml:"name"`
	LogLevel  string `yaml:"log_level"`
	ServerURL string `yaml:"server_url"`
}

type SecurityConfig struct {
	TLS TLSConfig `yaml:"tls"`
}

type TLSConfig struct {
	CAFile             string `yaml:"ca_file"`
	CertFile           string `yaml:"cert_file"`
	KeyFile            string `yaml:"key_file"`
	InsecureSkipVerify bool   `yaml:"insecure_skip_verify"`
}

type PipelineConfig struct {
	Exporters ExportersMap `yaml:"exporters"`
	Routes    []RouteRule  `yaml:"routes"` // Future: Advanced routing
}

type ExportersMap struct {
	Splunk  SplunkConfig  `yaml:"splunk"`
	Wazuh   WazuhConfig   `yaml:"wazuh"`
	Webhook WebhookConfig `yaml:"webhook"`
}

type SplunkConfig struct {
	Enabled       bool          `yaml:"enabled"`
	URL           string        `yaml:"url"`
	Token         string        `yaml:"token"`
	BatchInterval time.Duration `yaml:"batch_interval"`
}

type WazuhConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Protocol string `yaml:"protocol"`
}

type WebhookConfig struct {
	Enabled  bool              `yaml:"enabled"`
	URL      string            `yaml:"url"`
	Method   string            `yaml:"method"`
	Headers  map[string]string `yaml:"headers"`
	Template string            `yaml:"template"`
}

type RouteRule struct {
	DataTypes    []string `yaml:"data_types"`
	Destinations []string `yaml:"destinations"`
}

type ModulesConfig struct {
	Osquery   OsqueryConfig   `yaml:"osquery"`
	Hive      HiveConfig      `yaml:"hive"`
	Container ContainerConfig `yaml:"container"`
}

type OsqueryConfig struct {
	Enabled    bool   `yaml:"enabled"`
	BinaryPath string `yaml:"binary_path"`
}

type HiveConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Port     int    `yaml:"port"`
	CacheDir string `yaml:"cache_dir"`
}

type ContainerConfig struct {
	Enabled    bool   `yaml:"enabled"`
	SocketPath string `yaml:"socket_path"`
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

	// Environment Variable Overrides (Cloud-Native / Mass Deployment)
	if envName := os.Getenv("NEXUS_AGENT_NAME"); envName != "" {
		cfg.Agent.Name = envName
	}
	if envSplunkURL := os.Getenv("NEXUS_SPLUNK_URL"); envSplunkURL != "" {
		cfg.Pipeline.Exporters.Splunk.URL = envSplunkURL
		cfg.Pipeline.Exporters.Splunk.Enabled = true
	}
	if envSplunkToken := os.Getenv("NEXUS_SPLUNK_TOKEN"); envSplunkToken != "" {
		cfg.Pipeline.Exporters.Splunk.Token = envSplunkToken
	}

	return &cfg, nil
}

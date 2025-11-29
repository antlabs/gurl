package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// BatchConfig represents the batch test configuration
type BatchConfig struct {
	Version  string          `yaml:"version" json:"version"`
	Tests    []BatchTest     `yaml:"tests" json:"tests"`
	Notifier *NotifierConfig `yaml:"notifier,omitempty" json:"notifier,omitempty"`
}

// NotifierConfig defines configuration for batch result notifications.
type NotifierConfig struct {
	Type          string `yaml:"type" json:"type"`
	Enabled       bool   `yaml:"enabled" json:"enabled"`
	OnlyOnFail    bool   `yaml:"only_on_fail" json:"only_on_fail"`
	FeishuWebhook string `yaml:"feishu_webhook" json:"feishu_webhook"`
	Title         string `yaml:"title" json:"title"`
}

// BatchTest represents a single test in the batch
type BatchTest struct {
	Name        string `yaml:"name" json:"name"`
	Curl        string `yaml:"curl" json:"curl"`
	Connections int    `yaml:"connections,omitempty" json:"connections,omitempty"`
	Duration    string `yaml:"duration,omitempty" json:"duration,omitempty"`
	Threads     int    `yaml:"threads,omitempty" json:"threads,omitempty"`
	Rate        int    `yaml:"rate,omitempty" json:"rate,omitempty"`
	Timeout     string `yaml:"timeout,omitempty" json:"timeout,omitempty"`
	Verbose     bool   `yaml:"verbose,omitempty" json:"verbose,omitempty"`
	UseNetHTTP  bool   `yaml:"use_nethttp,omitempty" json:"use_nethttp,omitempty"`
	Asserts     string `yaml:"asserts,omitempty" json:"asserts,omitempty"`
	Requests    int64  `yaml:"requests,omitempty" json:"requests,omitempty"`
}

// ToConfig converts BatchTest to Config with defaults
func (bt *BatchTest) ToConfig(defaults *Config) (*Config, error) {
	cfg := &Config{
		// Copy defaults
		Connections:  defaults.Connections,
		Duration:     defaults.Duration,
		Threads:      defaults.Threads,
		Rate:         defaults.Rate,
		Timeout:      defaults.Timeout,
		Verbose:      defaults.Verbose,
		UseNetHTTP:   defaults.UseNetHTTP,
		PrintLatency: defaults.PrintLatency,
		Requests:     defaults.Requests,
	}

	if bt.Requests > 0 {
		cfg.Requests = bt.Requests
	}
	// Override with batch test specific values
	if bt.Connections > 0 {
		cfg.Connections = bt.Connections
	}
	if bt.Duration != "" {
		duration, err := time.ParseDuration(bt.Duration)
		if err != nil {
			return nil, fmt.Errorf("invalid duration '%s' for test '%s': %v", bt.Duration, bt.Name, err)
		}
		cfg.Duration = duration
	}
	if bt.Threads > 0 {
		cfg.Threads = bt.Threads
	}
	if bt.Rate > 0 {
		cfg.Rate = bt.Rate
	}
	if bt.Timeout != "" {
		timeout, err := time.ParseDuration(bt.Timeout)
		if err != nil {
			return nil, fmt.Errorf("invalid timeout '%s' for test '%s': %v", bt.Timeout, bt.Name, err)
		}
		cfg.Timeout = timeout
	}
	if bt.Verbose {
		cfg.Verbose = bt.Verbose
	}
	if bt.UseNetHTTP {
		cfg.UseNetHTTP = bt.UseNetHTTP
	}

	// Set curl command
	cfg.CurlCommand = bt.Curl

	// Set asserts text for this test (if any)
	cfg.Asserts = bt.Asserts

	return cfg, nil
}

// LoadBatchConfig loads batch configuration from file
func LoadBatchConfig(filename string) (*BatchConfig, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config BatchConfig
	ext := strings.ToLower(filepath.Ext(filename))

	switch ext {
	case ".yaml", ".yml":
		err = yaml.Unmarshal(data, &config)
	case ".json":
		err = json.Unmarshal(data, &config)
	default:
		return nil, fmt.Errorf("unsupported config file format: %s (supported: .yaml, .yml, .json)", ext)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	return &config, nil
}

// Validate validates the batch configuration
func (bc *BatchConfig) Validate() error {
	if bc.Version == "" {
		return fmt.Errorf("version is required")
	}

	if len(bc.Tests) == 0 {
		return fmt.Errorf("at least one test is required")
	}

	for i, test := range bc.Tests {
		if test.Name == "" {
			return fmt.Errorf("test[%d]: name is required", i)
		}
		if test.Curl == "" {
			return fmt.Errorf("test[%d] (%s): curl command is required", i, test.Name)
		}
		if test.Connections < 0 {
			return fmt.Errorf("test[%d] (%s): connections cannot be negative", i, test.Name)
		}
		if test.Threads < 0 {
			return fmt.Errorf("test[%d] (%s): threads cannot be negative", i, test.Name)
		}
		if test.Rate < 0 {
			return fmt.Errorf("test[%d] (%s): rate cannot be negative", i, test.Name)
		}
		// Validate duration format if provided
		if test.Duration != "" {
			if _, err := time.ParseDuration(test.Duration); err != nil {
				return fmt.Errorf("test[%d] (%s): invalid duration format '%s'", i, test.Name, test.Duration)
			}
		}
		// Validate timeout format if provided
		if test.Timeout != "" {
			if _, err := time.ParseDuration(test.Timeout); err != nil {
				return fmt.Errorf("test[%d] (%s): invalid timeout format '%s'", i, test.Name, test.Timeout)
			}
		}
	}

	return nil
}

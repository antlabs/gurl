package mock

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// MockConfig represents the mock server configuration file
type MockConfig struct {
	Port   int           `json:"port" yaml:"port"`
	Routes []RouteConfig `json:"routes" yaml:"routes"`
}

// LoadConfig loads mock server configuration from a file
func LoadConfig(filePath string) (*MockConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	config := &MockConfig{}
	ext := filepath.Ext(filePath)

	switch ext {
	case ".json":
		if err := json.Unmarshal(data, config); err != nil {
			return nil, fmt.Errorf("failed to parse JSON config: %w", err)
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, config); err != nil {
			return nil, fmt.Errorf("failed to parse YAML config: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported config file format: %s (use .json or .yaml)", ext)
	}

	return config, nil
}

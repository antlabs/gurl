package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// CompareRequest represents a single request definition used in compare scenarios.
type CompareRequest struct {
	Name  string `yaml:"name" json:"name"`
	Curl  string `yaml:"curl" json:"curl"`
	Group string `yaml:"group,omitempty" json:"group,omitempty"`
	Role  string `yaml:"role,omitempty" json:"role,omitempty"`
}

// RequestSet represents a named set of requests used for pair_by_index mode.
type RequestSet struct {
	Name     string           `yaml:"name" json:"name"`
	Requests []CompareRequest `yaml:"requests" json:"requests"`
}

// CompareScenario describes a single compare scenario.
type CompareScenario struct {
	Name            string   `yaml:"name" json:"name"`
	Mode            string   `yaml:"mode,omitempty" json:"mode,omitempty"`

	// one_to_one
	Base   string `yaml:"base,omitempty" json:"base,omitempty"`
	Target string `yaml:"target,omitempty" json:"target,omitempty"`

	// one_to_many
	Targets []string `yaml:"targets,omitempty" json:"targets,omitempty"`

	// pair_by_index
	LeftSet  string `yaml:"left_set,omitempty" json:"left_set,omitempty"`
	RightSet string `yaml:"right_set,omitempty" json:"right_set,omitempty"`

	// group_by_field
	GroupField string `yaml:"group_field,omitempty" json:"group_field,omitempty"`
	BaseRole   string `yaml:"base_role,omitempty" json:"base_role,omitempty"`
	TargetRole string `yaml:"target_role,omitempty" json:"target_role,omitempty"`

	ResponseCompare string `yaml:"response_compare" json:"response_compare"`
}

// CompareConfig is the top-level configuration for compare mode.
type CompareConfig struct {
	Requests    []CompareRequest        `yaml:"requests,omitempty" json:"requests,omitempty"`
	RequestSets map[string][]CompareRequest `yaml:"request_sets,omitempty" json:"request_sets,omitempty"`
	Scenarios   []CompareScenario       `yaml:"compare" json:"compare"`
}

// LoadCompareConfig loads compare configuration from a YAML or JSON file.
func LoadCompareConfig(filename string) (*CompareConfig, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read compare config file: %v", err)
	}

	var cfg CompareConfig
	ext := strings.ToLower(filepath.Ext(filename))

	switch ext {
	case ".yaml", ".yml":
		err = yaml.Unmarshal(data, &cfg)
	case ".json":
		err = json.Unmarshal(data, &cfg)
	default:
		return nil, fmt.Errorf("unsupported compare config file format: %s (supported: .yaml, .yml, .json)", ext)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to parse compare config file: %v", err)
	}

	return &cfg, nil
}

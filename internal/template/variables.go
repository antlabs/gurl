package template

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

// VariableGenerator generates dynamic values for template variables
type VariableGenerator struct {
	sequences map[string]*int64 // For sequence variables
}

// NewVariableGenerator creates a new variable generator
func NewVariableGenerator() *VariableGenerator {
	return &VariableGenerator{
		sequences: make(map[string]*int64),
	}
}

// GenerateValue generates a value based on the variable type and parameters
func (vg *VariableGenerator) GenerateValue(varType, params string) (string, error) {
	switch varType {
	case "random":
		return vg.generateRandom(params)
	case "uuid":
		return vg.generateUUID()
	case "timestamp":
		return vg.generateTimestamp(params)
	case "sequence":
		return vg.generateSequence(params)
	case "choice":
		return vg.generateChoice(params)
	case "now":
		return vg.generateNow(params)
	default:
		return "", fmt.Errorf("unknown variable type: %s", varType)
	}
}

// generateRandom generates a random number within the specified range
// Format: "min-max" or just "max" (min defaults to 0)
func (vg *VariableGenerator) generateRandom(params string) (string, error) {
	if params == "" {
		// Default range 1-1000
		params = "1-1000"
	}

	var min, max int64
	var err error

	if strings.Contains(params, "-") {
		parts := strings.Split(params, "-")
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid random format: %s (expected min-max)", params)
		}
		min, err = strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return "", fmt.Errorf("invalid min value: %s", parts[0])
		}
		max, err = strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return "", fmt.Errorf("invalid max value: %s", parts[1])
		}
	} else {
		min = 0
		max, err = strconv.ParseInt(params, 10, 64)
		if err != nil {
			return "", fmt.Errorf("invalid max value: %s", params)
		}
	}

	if min >= max {
		return "", fmt.Errorf("min value (%d) must be less than max value (%d)", min, max)
	}

	// Generate random number in range [min, max]
	rangeSize := max - min + 1
	n, err := rand.Int(rand.Reader, big.NewInt(rangeSize))
	if err != nil {
		return "", fmt.Errorf("failed to generate random number: %v", err)
	}

	result := min + n.Int64()
	return strconv.FormatInt(result, 10), nil
}

// generateUUID generates a new UUID
func (vg *VariableGenerator) generateUUID() (string, error) {
	u := uuid.New()
	return u.String(), nil
}

// generateTimestamp generates a timestamp
// Format: "unix" for unix timestamp, "rfc3339" for RFC3339 format, or empty for unix
func (vg *VariableGenerator) generateTimestamp(params string) (string, error) {
	now := time.Now()
	
	switch params {
	case "", "unix":
		return strconv.FormatInt(now.Unix(), 10), nil
	case "unix_ms":
		return strconv.FormatInt(now.UnixMilli(), 10), nil
	case "unix_ns":
		return strconv.FormatInt(now.UnixNano(), 10), nil
	case "rfc3339":
		return now.Format(time.RFC3339), nil
	case "iso8601":
		return now.Format("2006-01-02T15:04:05Z07:00"), nil
	case "date":
		return now.Format("2006-01-02"), nil
	case "time":
		return now.Format("15:04:05"), nil
	default:
		// Try to parse as custom format
		return now.Format(params), nil
	}
}

// generateNow is an alias for generateTimestamp for backward compatibility
func (vg *VariableGenerator) generateNow(params string) (string, error) {
	if params == "" {
		params = "unix"
	}
	return vg.generateTimestamp(params)
}

// generateSequence generates an incrementing sequence number
// Format: "start" where start is the initial value (default 1)
func (vg *VariableGenerator) generateSequence(params string) (string, error) {
	start := int64(1)
	if params != "" {
		var err error
		start, err = strconv.ParseInt(params, 10, 64)
		if err != nil {
			return "", fmt.Errorf("invalid sequence start value: %s", params)
		}
	}

	// Use params as key to support multiple sequences
	key := fmt.Sprintf("seq_%s", params)
	if _, exists := vg.sequences[key]; !exists {
		vg.sequences[key] = &start
		return strconv.FormatInt(start, 10), nil
	}

	// Atomically increment and return
	next := atomic.AddInt64(vg.sequences[key], 1)
	return strconv.FormatInt(next, 10), nil
}

// generateChoice randomly selects from a comma-separated list of options
// Format: "opt1,opt2,opt3"
func (vg *VariableGenerator) generateChoice(params string) (string, error) {
	if params == "" {
		return "", fmt.Errorf("choice variable requires options")
	}

	options := strings.Split(params, ",")
	if len(options) == 0 {
		return "", fmt.Errorf("choice variable requires at least one option")
	}

	// Trim whitespace from options
	for i, opt := range options {
		options[i] = strings.TrimSpace(opt)
	}

	// Generate random index
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(options))))
	if err != nil {
		return "", fmt.Errorf("failed to generate random choice: %v", err)
	}

	return options[n.Int64()], nil
}

// PredefinedVariables contains commonly used variable definitions
var PredefinedVariables = map[string]string{
	"user_id":    "random:1-10000",
	"session_id": "uuid",
	"timestamp":  "now:unix",
	"page":       "random:1-100",
	"limit":      "choice:10,20,50,100",
}

// VariableContext holds variable definitions and their generators
type VariableContext struct {
	generator *VariableGenerator
	variables map[string]string // variable name -> definition
}

// NewVariableContext creates a new variable context
func NewVariableContext() *VariableContext {
	return &VariableContext{
		generator: NewVariableGenerator(),
		variables: make(map[string]string),
	}
}

// SetVariable sets a variable definition
func (vc *VariableContext) SetVariable(name, definition string) {
	vc.variables[name] = definition
}

// GetVariable gets a variable definition
func (vc *VariableContext) GetVariable(name string) (string, bool) {
	def, exists := vc.variables[name]
	return def, exists
}

// GenerateVariableValue generates a value for a named variable
func (vc *VariableContext) GenerateVariableValue(name string) (string, error) {
	definition, exists := vc.variables[name]
	if !exists {
		// Check predefined variables
		if predefined, ok := PredefinedVariables[name]; ok {
			definition = predefined
		} else {
			return "", fmt.Errorf("undefined variable: %s", name)
		}
	}

	// Parse definition (format: "type:params")
	parts := strings.SplitN(definition, ":", 2)
	varType := parts[0]
	params := ""
	if len(parts) > 1 {
		params = parts[1]
	}

	return vc.generator.GenerateValue(varType, params)
}

// ListVariables returns all defined variables
func (vc *VariableContext) ListVariables() map[string]string {
	result := make(map[string]string)
	
	// Add user-defined variables
	for name, def := range vc.variables {
		result[name] = def
	}
	
	// Add predefined variables that aren't overridden
	for name, def := range PredefinedVariables {
		if _, exists := result[name]; !exists {
			result[name] = def
		}
	}
	
	return result
}

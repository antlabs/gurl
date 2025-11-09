package template

import (
	"fmt"
	"regexp"
	"strings"
)

// TemplateParser handles parsing and processing of template strings
type TemplateParser struct {
	context *VariableContext
}

// NewTemplateParser creates a new template parser
func NewTemplateParser() *TemplateParser {
	return &TemplateParser{
		context: NewVariableContext(),
	}
}

// NewTemplateParserWithContext creates a new template parser with existing context
func NewTemplateParserWithContext(context *VariableContext) *TemplateParser {
	return &TemplateParser{
		context: context,
	}
}

// SetVariable sets a variable definition in the parser context
func (tp *TemplateParser) SetVariable(name, definition string) {
	tp.context.SetVariable(name, definition)
}

// GetContext returns the variable context
func (tp *TemplateParser) GetContext() *VariableContext {
	return tp.context
}

// Template variable patterns
var (
	// Combined pattern for all template variables
	// Matches {{.variable}}, {{variable}}, {{.function:params}}, or {{function:params}}
	templatePattern = regexp.MustCompile(`\{\{\.?([a-zA-Z_][a-zA-Z0-9_]*)(?::([^}]*))?\}\}`)
)

// ParseTemplate processes a template string and replaces variables with generated values
func (tp *TemplateParser) ParseTemplate(template string) (string, error) {
	if template == "" {
		return template, nil
	}

	// Find all template variables
	matches := templatePattern.FindAllStringSubmatch(template, -1)
	if len(matches) == 0 {
		return template, nil
	}

	result := template

	// Process each match
	for _, match := range matches {
		fullMatch := match[0] // {{.variable}} or {{.function:params}}
		varName := match[1]   // variable or function name
		params := ""
		if len(match) > 2 {
			params = match[2] // parameters (if any)
		}

		// Generate value for this variable
		value, err := tp.generateValue(varName, params)
		if err != nil {
			return "", fmt.Errorf("failed to generate value for variable '%s': %v", varName, err)
		}

		// Replace the template variable with the generated value
		result = strings.ReplaceAll(result, fullMatch, value)
	}

	return result, nil
}

// generateValue generates a value for a variable, either from context or as a built-in function
func (tp *TemplateParser) generateValue(name, params string) (string, error) {
	// First, try to get from variable context
	if definition, exists := tp.context.GetVariable(name); exists {
		// If params are provided in template, they override the definition params
		if params != "" {
			// Parse the definition to get the type
			parts := strings.SplitN(definition, ":", 2)
			varType := parts[0]
			return tp.context.generator.GenerateValue(varType, params)
		}
		return tp.context.GenerateVariableValue(name)
	}

	// If not found in context, try as a built-in function
	return tp.context.generator.GenerateValue(name, params)
}

// ParseTemplateMultiple processes a template string multiple times and returns unique results
func (tp *TemplateParser) ParseTemplateMultiple(template string, count int) ([]string, error) {
	if count <= 0 {
		return nil, fmt.Errorf("count must be greater than 0")
	}

	results := make([]string, count)
	for i := 0; i < count; i++ {
		result, err := tp.ParseTemplate(template)
		if err != nil {
			return nil, fmt.Errorf("failed to parse template (iteration %d): %v", i+1, err)
		}
		results[i] = result
	}

	return results, nil
}

// ValidateTemplate checks if a template string is valid without generating values
func (tp *TemplateParser) ValidateTemplate(template string) error {
	if template == "" {
		return nil
	}

	matches := templatePattern.FindAllStringSubmatch(template, -1)

	for _, match := range matches {
		varName := match[1]
		params := ""
		if len(match) > 2 {
			params = match[2]
		}

		// Check if variable exists in context
		if _, exists := tp.context.GetVariable(varName); exists {
			continue
		}

		// Check if it's a predefined variable
		if _, exists := PredefinedVariables[varName]; exists {
			continue
		}

		// Check if it's a built-in function
		switch varName {
		case "random", "uuid", "timestamp", "sequence", "choice", "now":
			// Validate parameters for built-in functions
			if err := tp.validateBuiltinParams(varName, params); err != nil {
				return fmt.Errorf("invalid parameters for '%s': %v", varName, err)
			}
		default:
			return fmt.Errorf("undefined variable or function: %s", varName)
		}
	}

	return nil
}

// validateBuiltinParams validates parameters for built-in functions
func (tp *TemplateParser) validateBuiltinParams(funcName, params string) error {
	switch funcName {
	case "random":
		if params == "" {
			return nil // Default range is valid
		}
		if strings.Contains(params, "-") {
			parts := strings.Split(params, "-")
			if len(parts) != 2 {
				return fmt.Errorf("expected format 'min-max'")
			}
		}
		// Additional validation could be added here

	case "sequence":
		if params != "" {
			// Validate that params is a number
			if _, err := fmt.Sscanf(params, "%d", new(int)); err != nil {
				return fmt.Errorf("sequence start value must be a number")
			}
		}

	case "choice":
		if params == "" {
			return fmt.Errorf("choice function requires options")
		}
		if !strings.Contains(params, ",") && strings.TrimSpace(params) == "" {
			return fmt.Errorf("choice function requires at least one option")
		}

	case "timestamp", "now":
		// Most timestamp formats are valid, so we don't validate strictly

	case "uuid":
		// UUID doesn't take parameters
		if params != "" {
			return fmt.Errorf("uuid function does not accept parameters")
		}
	}

	return nil
}

// ExtractVariables extracts all variable names from a template string
func (tp *TemplateParser) ExtractVariables(template string) []string {
	matches := templatePattern.FindAllStringSubmatch(template, -1)

	var variables []string
	seen := make(map[string]bool)

	for _, match := range matches {
		varName := match[1]
		if !seen[varName] {
			variables = append(variables, varName)
			seen[varName] = true
		}
	}

	return variables
}

// ReplaceVariables replaces template variables with provided values
func (tp *TemplateParser) ReplaceVariables(template string, values map[string]string) string {
	result := template

	for varName, value := range values {
		// Replace both {{.variable}} and {{variable}} patterns
		patterns := []string{
			fmt.Sprintf("{{.%s}}", varName),
			fmt.Sprintf("{{%s}}", varName),
		}

		for _, pattern := range patterns {
			result = strings.ReplaceAll(result, pattern, value)
		}
	}

	return result
}

// ParseVariableDefinition parses a variable definition string
// Format: "name=type:params" or "name=value"
func ParseVariableDefinition(definition string) (name, varType, params string, err error) {
	parts := strings.SplitN(definition, "=", 2)
	if len(parts) != 2 {
		return "", "", "", fmt.Errorf("invalid variable definition format: %s (expected name=definition)", definition)
	}

	name = strings.TrimSpace(parts[0])
	if name == "" {
		return "", "", "", fmt.Errorf("variable name cannot be empty")
	}

	defParts := strings.SplitN(parts[1], ":", 2)
	varType = strings.TrimSpace(defParts[0])
	if len(defParts) > 1 {
		params = strings.TrimSpace(defParts[1])
	}

	return name, varType, params, nil
}

// ParseVariableDefinitions parses multiple variable definitions
func ParseVariableDefinitions(definitions []string) (*VariableContext, error) {
	context := NewVariableContext()

	for _, def := range definitions {
		name, varType, params, err := ParseVariableDefinition(def)
		if err != nil {
			return nil, fmt.Errorf("failed to parse variable definition '%s': %v", def, err)
		}

		// Construct the full definition
		fullDef := varType
		if params != "" {
			fullDef = fmt.Sprintf("%s:%s", varType, params)
		}

		context.SetVariable(name, fullDef)
	}

	return context, nil
}

// HasTemplateVariables checks if a string contains template variables
func HasTemplateVariables(text string) bool {
	return templatePattern.MatchString(text)
}

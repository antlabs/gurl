package template

import (
	"fmt"
	"strings"
)

// BuiltinFunction represents a built-in template function
type BuiltinFunction struct {
	Name        string
	Description string
	Usage       string
	Examples    []string
}

// GetBuiltinFunctions returns a list of all built-in template functions
func GetBuiltinFunctions() []BuiltinFunction {
	return []BuiltinFunction{
		{
			Name:        "random",
			Description: "Generate a random number within specified range",
			Usage:       "{{random:min-max}} or {{random:max}}",
			Examples: []string{
				"{{random:1-100}}   // Random number between 1 and 100",
				"{{random:1000}}    // Random number between 0 and 1000",
				"{{random}}         // Random number between 1 and 1000 (default)",
			},
		},
		{
			Name:        "uuid",
			Description: "Generate a UUID (Universally Unique Identifier)",
			Usage:       "{{uuid}}",
			Examples: []string{
				"{{uuid}}           // e.g., 550e8400-e29b-41d4-a716-446655440000",
			},
		},
		{
			Name:        "timestamp",
			Description: "Generate a timestamp in various formats",
			Usage:       "{{timestamp:format}}",
			Examples: []string{
				"{{timestamp}}      // Unix timestamp (default)",
				"{{timestamp:unix}} // Unix timestamp: 1640995200",
				"{{timestamp:unix_ms}} // Unix timestamp in milliseconds",
				"{{timestamp:rfc3339}} // RFC3339 format: 2022-01-01T00:00:00Z",
				"{{timestamp:iso8601}} // ISO8601 format",
				"{{timestamp:date}}    // Date only: 2022-01-01",
				"{{timestamp:time}}    // Time only: 15:04:05",
			},
		},
		{
			Name:        "now",
			Description: "Alias for timestamp function",
			Usage:       "{{now:format}}",
			Examples: []string{
				"{{now}}            // Same as {{timestamp}}",
				"{{now:unix}}       // Unix timestamp",
			},
		},
		{
			Name:        "sequence",
			Description: "Generate incrementing sequence numbers",
			Usage:       "{{sequence:start}}",
			Examples: []string{
				"{{sequence}}       // Starts from 1: 1, 2, 3, ...",
				"{{sequence:100}}   // Starts from 100: 100, 101, 102, ...",
			},
		},
		{
			Name:        "choice",
			Description: "Randomly select from a list of options",
			Usage:       "{{choice:option1,option2,option3}}",
			Examples: []string{
				"{{choice:GET,POST,PUT}} // Randomly selects GET, POST, or PUT",
				"{{choice:apple,banana,orange}} // Randomly selects a fruit",
			},
		},
	}
}

// GetFunctionHelp returns help text for a specific function
func GetFunctionHelp(name string) (string, error) {
	functions := GetBuiltinFunctions()

	for _, fn := range functions {
		if fn.Name == name {
			var help strings.Builder
			help.WriteString(fmt.Sprintf("Function: %s\n", fn.Name))
			help.WriteString(fmt.Sprintf("Description: %s\n", fn.Description))
			help.WriteString(fmt.Sprintf("Usage: %s\n", fn.Usage))
			help.WriteString("Examples:\n")
			for _, example := range fn.Examples {
				help.WriteString(fmt.Sprintf("  %s\n", example))
			}
			return help.String(), nil
		}
	}

	return "", fmt.Errorf("unknown function: %s", name)
}

// GetAllFunctionsHelp returns help text for all functions
func GetAllFunctionsHelp() string {
	var help strings.Builder
	help.WriteString("Built-in Template Functions:\n")
	help.WriteString("===========================\n\n")

	functions := GetBuiltinFunctions()
	for i, fn := range functions {
		if i > 0 {
			help.WriteString("\n")
		}

		help.WriteString(fmt.Sprintf("%d. %s\n", i+1, fn.Name))
		help.WriteString(fmt.Sprintf("   Description: %s\n", fn.Description))
		help.WriteString(fmt.Sprintf("   Usage: %s\n", fn.Usage))
		help.WriteString("   Examples:\n")
		for _, example := range fn.Examples {
			help.WriteString(fmt.Sprintf("     %s\n", example))
		}
	}

	help.WriteString("\nTemplate Syntax:\n")
	help.WriteString("================\n")
	help.WriteString("- Use {{function}} or {{.function}} for functions without parameters\n")
	help.WriteString("- Use {{function:params}} or {{.function:params}} for functions with parameters\n")
	help.WriteString("- Variables can be defined with --var name=type:params\n")
	help.WriteString("- Variables can be referenced as {{name}} or {{.name}}\n")

	return help.String()
}

// ValidateFunction checks if a function name is valid and supported
func ValidateFunction(name string) bool {
	functions := GetBuiltinFunctions()
	for _, fn := range functions {
		if fn.Name == name {
			return true
		}
	}
	return false
}

// GetFunctionNames returns a list of all built-in function names
func GetFunctionNames() []string {
	functions := GetBuiltinFunctions()
	names := make([]string, len(functions))
	for i, fn := range functions {
		names[i] = fn.Name
	}
	return names
}

// TemplateExample represents a complete template usage example
type TemplateExample struct {
	Name        string
	Description string
	Template    string
	Variables   map[string]string
	Expected    string
}

// GetTemplateExamples returns example templates with their usage
func GetTemplateExamples() []TemplateExample {
	return []TemplateExample{
		{
			Name:        "User API with Random ID",
			Description: "Access user API with random user ID",
			Template:    "https://api.example.com/users/{{random:1-10000}}",
			Variables:   nil,
			Expected:    "https://api.example.com/users/7423",
		},
		{
			Name:        "Session-based Request",
			Description: "Request with session ID and timestamp",
			Template:    "https://api.example.com/data?session={{uuid}}&timestamp={{timestamp:unix}}",
			Variables:   nil,
			Expected:    "https://api.example.com/data?session=550e8400-e29b-41d4-a716-446655440000&timestamp=1640995200",
		},
		{
			Name:        "Pagination with Sequence",
			Description: "Paginated requests with incrementing page numbers",
			Template:    "https://api.example.com/items?page={{sequence:1}}&limit=20",
			Variables:   nil,
			Expected:    "https://api.example.com/items?page=1&limit=20 (then page=2, page=3, etc.)",
		},
		{
			Name:        "Random HTTP Method",
			Description: "Randomly choose HTTP method for testing",
			Template:    "curl -X {{choice:GET,POST,PUT,DELETE}} https://api.example.com/resource/{{random:1-100}}",
			Variables:   nil,
			Expected:    "curl -X POST https://api.example.com/resource/42",
		},
		{
			Name:        "Custom Variables",
			Description: "Using custom defined variables",
			Template:    "https://api.example.com/{{endpoint}}/{{user_id}}?token={{session_token}}",
			Variables: map[string]string{
				"endpoint":      "choice:users,orders,products",
				"user_id":       "random:1-1000",
				"session_token": "uuid",
			},
			Expected: "https://api.example.com/users/123?token=550e8400-e29b-41d4-a716-446655440000",
		},
		{
			Name:        "E-commerce Simulation",
			Description: "Simulate e-commerce API calls with realistic data",
			Template:    "curl -X POST https://shop.example.com/api/orders -d '{\"user_id\":{{user_id}},\"product_id\":{{product_id}},\"quantity\":{{quantity}},\"timestamp\":\"{{timestamp:rfc3339}}\"}'",
			Variables: map[string]string{
				"user_id":    "random:1-10000",
				"product_id": "random:100-999",
				"quantity":   "choice:1,2,3,4,5",
			},
			Expected: `curl -X POST https://shop.example.com/api/orders -d '{"user_id":7423,"product_id":456,"quantity":2,"timestamp":"2022-01-01T15:04:05Z"}'`,
		},
	}
}

// PrintTemplateExamples prints formatted template examples
func PrintTemplateExamples() string {
	var output strings.Builder
	output.WriteString("Template Usage Examples:\n")
	output.WriteString("========================\n\n")

	examples := GetTemplateExamples()
	for i, example := range examples {
		output.WriteString(fmt.Sprintf("%d. %s\n", i+1, example.Name))
		output.WriteString(fmt.Sprintf("   Description: %s\n", example.Description))
		output.WriteString(fmt.Sprintf("   Template: %s\n", example.Template))

		if len(example.Variables) > 0 {
			output.WriteString("   Variables:\n")
			for name, def := range example.Variables {
				output.WriteString(fmt.Sprintf("     --var %s=%s\n", name, def))
			}
		}

		output.WriteString(fmt.Sprintf("   Example Output: %s\n", example.Expected))

		if i < len(examples)-1 {
			output.WriteString("\n")
		}
	}

	return output.String()
}

// GetQuickStartGuide returns a quick start guide for template variables
func GetQuickStartGuide() string {
	return `Template Variables Quick Start Guide:
=====================================

1. Basic Usage:
   gurl -c 10 -d 30s 'https://api.example.com/users/{{random:1-1000}}'

2. With Custom Variables:
   gurl --var user_id=random:1-10000 --var method=choice:GET,POST 'https://api.example.com/{{method}}/users/{{user_id}}'

3. Multiple Variables:
   gurl --var session=uuid --var timestamp=now:unix 'https://api.example.com/data?session={{session}}&t={{timestamp}}'

4. In Batch Configuration:
   version: "1.0"
   tests:
     - name: "Dynamic User Test"
       curl: 'curl https://api.example.com/users/{{random:1-1000}}'
       connections: 50
       duration: "30s"

5. Available Functions:
   - {{random:1-100}}     # Random number 1-100
   - {{uuid}}             # UUID
   - {{timestamp:unix}}   # Unix timestamp
   - {{sequence:1}}       # Incrementing: 1,2,3...
   - {{choice:a,b,c}}     # Random choice from a,b,c

For more help: gurl --help-templates
`
}

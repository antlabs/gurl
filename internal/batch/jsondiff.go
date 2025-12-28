package batch

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/wI2L/jsondiff"
)

// JSONDiffResult represents the result of a JSON comparison
type JSONDiffResult struct {
	PairName    string
	BaseName    string
	TargetName  string
	Passed      bool
	Message     string
	Differences []FieldDifference
}

// FieldDifference represents a difference between two JSON fields
type FieldDifference struct {
	Field       string
	BaseValue   string
	TargetValue string
}

// CompareJSONResponses compares two JSON responses based on the configuration
func CompareJSONResponses(baseName, targetName string, baseBody, targetBody []byte, compareField string, ignoreFields []string) JSONDiffResult {
	return CompareJSONResponsesWithFields(baseName, targetName, baseBody, targetBody, compareField, compareField, ignoreFields)
}

// CompareJSONResponsesWithFields compares two JSON responses with separate field paths for base and target
func CompareJSONResponsesWithFields(baseName, targetName string, baseBody, targetBody []byte, baseField, targetField string, ignoreFields []string) JSONDiffResult {
	result := JSONDiffResult{
		BaseName:    baseName,
		TargetName:  targetName,
		Differences: make([]FieldDifference, 0),
	}

	// Extract fields from base and target responses
	var baseJSON, targetJSON []byte

	// Extract base field
	if baseField != "" {
		baseResult := gjson.GetBytes(baseBody, baseField)
		if !baseResult.Exists() {
			result.Passed = false
			result.Message = fmt.Sprintf("field '%s' not found in base response", baseField)
			return result
		}
		baseJSON = []byte(baseResult.Raw)
	} else {
		baseJSON = baseBody
	}

	// Extract target field
	if targetField != "" {
		targetResult := gjson.GetBytes(targetBody, targetField)
		if !targetResult.Exists() {
			result.Passed = false
			result.Message = fmt.Sprintf("field '%s' not found in target response", targetField)
			return result
		}
		targetJSON = []byte(targetResult.Raw)
	} else {
		targetJSON = targetBody
	}

	// Use jsondiff library to compare
	patch, err := jsondiff.CompareJSON(baseJSON, targetJSON)
	if err != nil {
		result.Passed = false
		result.Message = fmt.Sprintf("failed to compare JSON: %v", err)
		return result
	}

	// Filter out ignored fields from patch operations
	filteredOps := filterIgnoredFields(patch, ignoreFields)

	// Convert patch operations to FieldDifference
	diffs := convertPatchToDifferences(filteredOps)
	result.Differences = diffs

	if len(diffs) == 0 {
		result.Passed = true
		if baseField != "" && targetField != "" && baseField != targetField {
			result.Message = fmt.Sprintf("fields '%s' (base) and '%s' (target) are identical", baseField, targetField)
		} else if baseField != "" {
			result.Message = fmt.Sprintf("field '%s' is identical in both responses", baseField)
		} else {
			result.Message = "responses are identical"
		}
	} else {
		result.Passed = false
		result.Message = fmt.Sprintf("found %d difference(s)", len(diffs))
	}

	return result
}

// filterIgnoredFields filters out patch operations for ignored fields
func filterIgnoredFields(patch jsondiff.Patch, ignoreFields []string) jsondiff.Patch {
	if len(ignoreFields) == 0 {
		return patch
	}

	filteredOps := make(jsondiff.Patch, 0)
	for _, op := range patch {
		path := op.Path
		// Convert JSON Pointer path to dot notation (remove leading / and replace / with .)
		dotPath := strings.TrimPrefix(path, "/")
		dotPath = strings.ReplaceAll(dotPath, "/", ".")

		if !shouldIgnoreField(dotPath, ignoreFields) {
			filteredOps = append(filteredOps, op)
		}
	}

	return filteredOps
}

// convertPatchToDifferences converts jsondiff patch operations to FieldDifference
func convertPatchToDifferences(patch jsondiff.Patch) []FieldDifference {
	diffs := make([]FieldDifference, 0)

	for _, op := range patch {
		path := op.Path
		// Convert JSON Pointer path to dot notation
		dotPath := strings.TrimPrefix(path, "/")
		dotPath = strings.ReplaceAll(dotPath, "/", ".")

		var baseValue, targetValue string

		switch op.Type {
		case "replace":
			baseValue = formatValue(op.OldValue)
			targetValue = formatValue(op.Value)
		case "add":
			baseValue = "<missing>"
			targetValue = formatValue(op.Value)
		case "remove":
			baseValue = formatValue(op.OldValue)
			targetValue = "<missing>"
		default:
			continue
		}

		diffs = append(diffs, FieldDifference{
			Field:       dotPath,
			BaseValue:   baseValue,
			TargetValue: targetValue,
		})
	}

	return diffs
}

// formatValue formats a value for display
func formatValue(v interface{}) string {
	if v == nil {
		return "null"
	}

	// Try to marshal as JSON for complex types
	if data, err := json.Marshal(v); err == nil {
		return string(data)
	}

	return fmt.Sprintf("%v", v)
}

// shouldIgnoreField checks if a field path should be ignored
func shouldIgnoreField(path string, ignoreFields []string) bool {
	if path == "" {
		return false
	}

	for _, ignorePattern := range ignoreFields {
		// Exact match
		if path == ignorePattern {
			return true
		}

		// Prefix match (for nested fields)
		if strings.HasPrefix(path, ignorePattern+".") {
			return true
		}

		// Suffix match (for field names in any object)
		if strings.HasSuffix(path, "."+ignorePattern) {
			return true
		}

		// Simple field name match (last part of path)
		parts := strings.Split(path, ".")
		if len(parts) > 0 && parts[len(parts)-1] == ignorePattern {
			return true
		}
	}

	return false
}

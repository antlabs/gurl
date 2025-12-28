package batch

import (
	"strings"
	"testing"
)

func TestCompareJSONResponses_Identical(t *testing.T) {
	base := []byte(`{"code":200,"msg":"success","data":{"user_id":123,"username":"test"}}`)
	target := []byte(`{"code":200,"msg":"success","data":{"user_id":123,"username":"test"}}`)

	result := CompareJSONResponses("base", "target", base, target, "", nil)

	if !result.Passed {
		t.Errorf("Expected comparison to pass, but got: %s", result.Message)
	}

	if len(result.Differences) != 0 {
		t.Errorf("Expected no differences, but got %d", len(result.Differences))
	}
}

func TestCompareJSONResponses_WithDifferences(t *testing.T) {
	base := []byte(`{"code":200,"msg":"success","data":{"user_id":123,"username":"test"}}`)
	target := []byte(`{"code":200,"msg":"success","data":{"user_id":456,"username":"test"}}`)

	result := CompareJSONResponses("base", "target", base, target, "", nil)

	if result.Passed {
		t.Errorf("Expected comparison to fail, but it passed")
	}

	if len(result.Differences) == 0 {
		t.Errorf("Expected differences, but got none")
	}

	// Check that user_id difference was detected
	found := false
	for _, diff := range result.Differences {
		if diff.Field == "data.user_id" {
			found = true
			if diff.BaseValue != "123" || diff.TargetValue != "456" {
				t.Errorf("Expected user_id difference: 123 vs 456, got: %s vs %s", diff.BaseValue, diff.TargetValue)
			}
		}
	}

	if !found {
		t.Errorf("Expected to find user_id difference")
	}
}

func TestCompareJSONResponses_CompareField(t *testing.T) {
	base := []byte(`{"code":200,"msg":"success","data":{"user_id":123,"username":"test"}}`)
	target := []byte(`{"code":200,"msg":"success","data":{"user_id":123,"username":"test"}}`)

	// Compare only the data field
	result := CompareJSONResponses("base", "target", base, target, "data", nil)

	if !result.Passed {
		t.Errorf("Expected comparison to pass, but got: %s", result.Message)
	}
}

func TestCompareJSONResponses_CompareFieldWithDifference(t *testing.T) {
	base := []byte(`{"code":200,"msg":"success","data":{"user_id":123,"username":"test"}}`)
	target := []byte(`{"code":200,"msg":"error","data":{"user_id":123,"username":"test"}}`)

	// Compare only the data field (should pass even though msg is different)
	result := CompareJSONResponses("base", "target", base, target, "data", nil)

	if !result.Passed {
		t.Errorf("Expected comparison to pass when comparing only data field, but got: %s", result.Message)
	}
}

func TestCompareJSONResponses_IgnoreFields(t *testing.T) {
	base := []byte(`{"code":200,"msg":"success","data":{"user_id":123,"username":"test","timestamp":"2024-01-01"}}`)
	target := []byte(`{"code":200,"msg":"success","data":{"user_id":123,"username":"test","timestamp":"2024-01-02"}}`)

	// Ignore timestamp field
	ignoreFields := []string{"timestamp"}
	result := CompareJSONResponses("base", "target", base, target, "", ignoreFields)

	if !result.Passed {
		t.Errorf("Expected comparison to pass when ignoring timestamp, but got: %s", result.Message)
	}
}

func TestCompareJSONResponses_IgnoreNestedFields(t *testing.T) {
	base := []byte(`{"code":200,"data":{"user_id":123,"timestamp":"2024-01-01"}}`)
	target := []byte(`{"code":200,"data":{"user_id":123,"timestamp":"2024-01-02"}}`)

	// Ignore nested timestamp field
	ignoreFields := []string{"data.timestamp"}
	result := CompareJSONResponses("base", "target", base, target, "", ignoreFields)

	if !result.Passed {
		t.Errorf("Expected comparison to pass when ignoring data.timestamp, but got: %s", result.Message)
	}
}

func TestCompareJSONResponses_FieldNotFound(t *testing.T) {
	base := []byte(`{"code":200,"msg":"success"}`)
	target := []byte(`{"code":200,"msg":"success"}`)

	// Try to compare a field that doesn't exist
	result := CompareJSONResponses("base", "target", base, target, "data", nil)

	if result.Passed {
		t.Errorf("Expected comparison to fail when field not found")
	}

	if result.Message == "" {
		t.Errorf("Expected error message when field not found")
	}
}

func TestCompareJSONResponses_ArrayComparison(t *testing.T) {
	base := []byte(`{"items":[{"id":1,"name":"a"},{"id":2,"name":"b"}]}`)
	target := []byte(`{"items":[{"id":1,"name":"a"},{"id":2,"name":"b"}]}`)

	result := CompareJSONResponses("base", "target", base, target, "", nil)

	if !result.Passed {
		t.Errorf("Expected array comparison to pass, but got: %s", result.Message)
	}
}

func TestCompareJSONResponses_ArrayLengthDifference(t *testing.T) {
	base := []byte(`{"items":[{"id":1},{"id":2}]}`)
	target := []byte(`{"items":[{"id":1}]}`)

	result := CompareJSONResponses("base", "target", base, target, "", nil)

	if result.Passed {
		t.Errorf("Expected array comparison to fail due to length difference")
	}

	// jsondiff library detects array differences as remove operations
	// Check that some difference was detected (the library will report the removed element)
	if len(result.Differences) == 0 {
		t.Errorf("Expected to find differences for array length change, but got none")
	}
}

func TestShouldIgnoreField(t *testing.T) {
	ignoreFields := []string{"timestamp", "data.created_at", "request_id"}

	tests := []struct {
		path     string
		expected bool
	}{
		{"timestamp", true},
		{"data.timestamp", true},
		{"user.timestamp", true},
		{"data.created_at", true},
		{"request_id", true},
		{"data.request_id", true},
		{"user.name", false},
		{"data.user_id", false},
	}

	for _, tt := range tests {
		result := shouldIgnoreField(tt.path, ignoreFields)
		if result != tt.expected {
			t.Errorf("shouldIgnoreField(%q) = %v, expected %v", tt.path, result, tt.expected)
		}
	}
}

func TestCompareJSONResponsesWithFields_DifferentPaths(t *testing.T) {
	// Old version uses "result", new version uses "data"
	base := []byte(`{"code":200,"msg":"success","result":{"user_id":123,"username":"test"}}`)
	target := []byte(`{"code":200,"msg":"success","data":{"user_id":123,"username":"test"}}`)

	// Compare old version's "result" with new version's "data"
	result := CompareJSONResponsesWithFields("base", "target", base, target, "result", "data", nil)

	if !result.Passed {
		t.Errorf("Expected comparison to pass when comparing result vs data, but got: %s", result.Message)
	}

	if len(result.Differences) != 0 {
		t.Errorf("Expected no differences, but got %d", len(result.Differences))
	}
}

func TestCompareJSONResponsesWithFields_DifferentPathsWithDifference(t *testing.T) {
	// Old version uses "result", new version uses "data", but values are different
	base := []byte(`{"code":200,"result":{"user_id":123,"username":"test"}}`)
	target := []byte(`{"code":200,"data":{"user_id":456,"username":"test"}}`)

	result := CompareJSONResponsesWithFields("base", "target", base, target, "result", "data", nil)

	if result.Passed {
		t.Errorf("Expected comparison to fail due to user_id difference")
	}

	if len(result.Differences) == 0 {
		t.Errorf("Expected differences, but got none")
	}

	// Check that user_id difference was detected
	found := false
	for _, diff := range result.Differences {
		if diff.Field == "user_id" {
			found = true
			if diff.BaseValue != "123" || diff.TargetValue != "456" {
				t.Errorf("Expected user_id difference: 123 vs 456, got: %s vs %s", diff.BaseValue, diff.TargetValue)
			}
		}
	}

	if !found {
		t.Errorf("Expected to find user_id difference")
	}
}

func TestCompareJSONResponsesWithFields_NestedPaths(t *testing.T) {
	// Compare nested paths
	base := []byte(`{"response":{"v1":{"data":{"region":"Beijing","code":"156"}}}}`)
	target := []byte(`{"response":{"v2":{"data":{"region":"Beijing","code":"156"}}}}`)

	result := CompareJSONResponsesWithFields("base", "target", base, target, "response.v1.data", "response.v2.data", nil)

	if !result.Passed {
		t.Errorf("Expected nested path comparison to pass, but got: %s", result.Message)
	}
}

func TestCompareJSONResponsesWithFields_BaseFieldNotFound(t *testing.T) {
	base := []byte(`{"code":200,"msg":"success"}`)
	target := []byte(`{"code":200,"data":{"user_id":123}}`)

	result := CompareJSONResponsesWithFields("base", "target", base, target, "result", "data", nil)

	if result.Passed {
		t.Errorf("Expected comparison to fail when base field not found")
	}

	if !strings.Contains(result.Message, "not found in base response") {
		t.Errorf("Expected error message about base field not found, got: %s", result.Message)
	}
}

func TestCompareJSONResponsesWithFields_TargetFieldNotFound(t *testing.T) {
	base := []byte(`{"code":200,"result":{"user_id":123}}`)
	target := []byte(`{"code":200,"msg":"success"}`)

	result := CompareJSONResponsesWithFields("base", "target", base, target, "result", "data", nil)

	if result.Passed {
		t.Errorf("Expected comparison to fail when target field not found")
	}

	if !strings.Contains(result.Message, "not found in target response") {
		t.Errorf("Expected error message about target field not found, got: %s", result.Message)
	}
}

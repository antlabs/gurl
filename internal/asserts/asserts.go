package asserts

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

// HTTPResponse represents a single HTTP response used for assertions.
type HTTPResponse struct {
	Status   int
	Headers  http.Header
	Body     []byte
	Duration time.Duration
}

// Evaluate evaluates all assertions defined in the multi-line assertsText
// against the given HTTP response.
func Evaluate(assertsText string, resp *HTTPResponse) error {
	lines := strings.Split(assertsText, "\n")
	for i, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if err := evalSingle(line, resp); err != nil {
			return fmt.Errorf("assertion failed at line %d: %s: %w", i+1, line, err)
		}
	}
	return nil
}

func evalSingle(line string, resp *HTTPResponse) error {
	// Very small hand-written parser for the supported grammar.
	// Targets: status | header "Name" | gjson "Path" | body | duration_ms

	// Handle targets that start with a known prefix.
	if strings.HasPrefix(line, "status ") {
		return evalStatus(strings.TrimSpace(line[len("status"):]), resp)
	}

	if strings.HasPrefix(line, "header ") {
		rest := strings.TrimSpace(line[len("header"):])
		return evalHeader(rest, resp)
	}

	if strings.HasPrefix(line, "gjson ") {
		rest := strings.TrimSpace(line[len("gjson"):])
		return evalGJSONTarget(rest, resp)
	}

	if strings.HasPrefix(line, "body ") {
		rest := strings.TrimSpace(line[len("body"):])
		return evalBody(rest, resp)
	}

	if strings.HasPrefix(line, "duration_ms ") {
		rest := strings.TrimSpace(line[len("duration_ms"):])
		return evalDuration(rest, resp)
	}

	return fmt.Errorf("unsupported assert target in line: %s", line)
}

func evalStatus(expr string, resp *HTTPResponse) error {
	op, expectedStr, err := splitOperator(expr)
	if err != nil {
		return err
	}

	expected, err := strconv.Atoi(expectedStr)
	if err != nil {
		return fmt.Errorf("invalid status value '%s'", expectedStr)
	}

	return compareNumbers(float64(resp.Status), op, float64(expected))
}

func evalHeader(expr string, resp *HTTPResponse) error {
	// Expect: "Name" <op> <expected> | "Name" exists | "Name" not_exists
	name, rest, err := parseQuoted(expr)
	if err != nil {
		return fmt.Errorf("invalid header target: %w", err)
	}

	op, expected, err := splitOperator(rest)
	if err != nil {
		// Might be exists/not_exists without value
		word := strings.TrimSpace(rest)
		lower := strings.ToLower(word)
		value := resp.Headers.Get(name)
		if lower == "exists" {
			if value == "" {
				return fmt.Errorf("header '%s' does not exist", name)
			}
			return nil
		}
		if lower == "not_exists" {
			if value == "" {
				return nil
			}
			return fmt.Errorf("header '%s' exists", name)
		}
		return err
	}

	actual := resp.Headers.Get(name)
	return compareStrings(actual, op, expected)
}

func evalBody(expr string, resp *HTTPResponse) error {
	op, expected, err := splitOperator(expr)
	if err != nil {
		return err
	}

	body := string(resp.Body)
	return compareStrings(body, op, expected)
}

func evalDuration(expr string, resp *HTTPResponse) error {
	op, expectedStr, err := splitOperator(expr)
	if err != nil {
		return err
	}

	expected, err := strconv.ParseFloat(expectedStr, 64)
	if err != nil {
		return fmt.Errorf("invalid duration_ms value '%s'", expectedStr)
	}

	actual := float64(resp.Duration.Milliseconds())
	return compareNumbers(actual, op, expected)
}

func evalGJSONTarget(expr string, resp *HTTPResponse) error {
	// Expect: "path" <op> <expected>
	path, rest, err := parseQuoted(expr)
	if err != nil {
		return fmt.Errorf("invalid gjson target: %w", err)
	}

	res := gjson.GetBytes(resp.Body, path)

	op, expectedToken, err := splitOperator(rest)
	if err != nil {
		// handle exists/not_exists
		word := strings.TrimSpace(rest)
		lower := strings.ToLower(word)
		if lower == "exists" {
			if !res.Exists() {
				return fmt.Errorf("gjson path '%s' does not exist", path)
			}
			return nil
		}
		if lower == "not_exists" {
			if !res.Exists() {
				return nil
			}
			return fmt.Errorf("gjson path '%s' exists", path)
		}
		return err
	}

	if !res.Exists() {
		return fmt.Errorf("gjson path '%s' does not exist", path)
	}

	// Decide comparison type based on expected literal
	if strings.HasPrefix(expectedToken, "/") && strings.HasSuffix(expectedToken, "/") && len(expectedToken) >= 2 {
		// regex
		pattern := expectedToken[1 : len(expectedToken)-1]
		re, err := regexp.Compile(pattern)
		if err != nil {
			return fmt.Errorf("invalid regex '%s'", pattern)
		}
		actualStr := res.String()
		if op != "matches" {
			return fmt.Errorf("operator '%s' not supported with regex", op)
		}
		if !re.MatchString(actualStr) {
			return fmt.Errorf("regex '%s' does not match '%s'", pattern, actualStr)
		}
		return nil
	}

	// boolean literal
	if expectedToken == "true" || expectedToken == "false" {
		bExpected := expectedToken == "true"
		bActual := res.Bool()
		return compareBools(bActual, op, bExpected)
	}

	// numeric literal
	if n, errNum := strconv.ParseFloat(expectedToken, 64); errNum == nil {
		actualNum := res.Float()
		return compareNumbers(actualNum, op, n)
	}

	// string literal (wrapped in quotes) or plain token
	str, err := unquoteMaybe(expectedToken)
	if err != nil {
		return err
	}

	actualStr := res.String()
	return compareStrings(actualStr, op, str)
}

// --- helpers ---

func splitOperator(expr string) (op string, rhs string, err error) {
	expr = strings.TrimSpace(expr)
	// try multi-char operators first
	for _, candidate := range []string{"==", "!=", ">=", "<=", ">", "<", "contains", "not_contains", "starts_with", "ends_with", "matches"} {
		if strings.HasPrefix(expr, candidate+" ") {
			return candidate, strings.TrimSpace(expr[len(candidate):]), nil
		}
	}
	return "", "", fmt.Errorf("unsupported or missing operator in '%s'", expr)
}

func compareNumbers(actual float64, op string, expected float64) error {
	switch op {
	case "==":
		if actual == expected {
			return nil
		}
	case "!=":
		if actual != expected {
			return nil
		}
	case ">":
		if actual > expected {
			return nil
		}
	case ">=":
		if actual >= expected {
			return nil
		}
	case "<":
		if actual < expected {
			return nil
		}
	case "<=":
		if actual <= expected {
			return nil
		}
	default:
		return fmt.Errorf("operator '%s' not supported for numeric comparison", op)
	}
	return fmt.Errorf("actual=%v, expected %s %v", actual, op, expected)
}

func compareBools(actual bool, op string, expected bool) error {
	switch op {
	case "==":
		if actual == expected {
			return nil
		}
	case "!=":
		if actual != expected {
			return nil
		}
	default:
		return fmt.Errorf("operator '%s' not supported for bool comparison", op)
	}
	return fmt.Errorf("actual=%v, expected %s %v", actual, op, expected)
}

func compareStrings(actual, op, expected string) error {
	switch op {
	case "==":
		if actual == expected {
			return nil
		}
	case "!=":
		if actual != expected {
			return nil
		}
	case "contains":
		if strings.Contains(actual, expected) {
			return nil
		}
	case "not_contains":
		if !strings.Contains(actual, expected) {
			return nil
		}
	case "starts_with":
		if strings.HasPrefix(actual, expected) {
			return nil
		}
	case "ends_with":
		if strings.HasSuffix(actual, expected) {
			return nil
		}
	default:
		return fmt.Errorf("operator '%s' not supported for string comparison", op)
	}
	return fmt.Errorf("actual='%s', expected %s '%s'", actual, op, expected)
}

func parseQuoted(s string) (value string, rest string, err error) {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, "\"") {
		return "", "", fmt.Errorf("expected quoted string, got '%s'", s)
	}
	// find closing quote
	idx := strings.Index(s[1:], "\"")
	if idx < 0 {
		return "", "", fmt.Errorf("unterminated quoted string in '%s'", s)
	}
	idx++ // account for starting offset
	value = s[1:idx]
	rest = strings.TrimSpace(s[idx+1:])
	return value, rest, nil
}

func unquoteMaybe(s string) (string, error) {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"") && len(s) >= 2 {
		unquoted, err := strconv.Unquote(s)
		if err != nil {
			return "", fmt.Errorf("invalid quoted string %s", s)
		}
		return unquoted, nil
	}
	return s, nil
}

func toFloat64(v any) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case float32:
		return float64(n), true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	case int32:
		return float64(n), true
	case uint:
		return float64(n), true
	case uint64:
		return float64(n), true
	default:
		return 0, false
	}
}

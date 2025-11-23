package asserts

import (
	"net/http"
	"testing"
	"time"
)

func TestEvaluateGJSONEquals(t *testing.T) {
	resp := &HTTPResponse{
		Status:   200,
		Headers:  http.Header{"Content-Type": []string{"application/json"}},
		Body:     []byte(`{"message":"Fast response","latency":"0ms"}`),
		Duration: 10 * time.Millisecond,
	}

	assertText := `
	gjson "message" == "Fast response"
	gjson "latency" == "0ms"
	`

	if err := Evaluate(assertText, resp); err != nil {
		t.Fatalf("expected assertions to pass, got error: %v", err)
	}
}

func TestEvaluateGJSONNotEquals(t *testing.T) {
	resp := &HTTPResponse{
		Status:   200,
		Headers:  http.Header{"Content-Type": []string{"application/json"}},
		Body:     []byte(`{"message":"Fast response","latency":"0ms"}`),
		Duration: 10 * time.Millisecond,
	}

	assertText := `
	gjson "message" != "Fast response"
	`

	if err := Evaluate(assertText, resp); err == nil {
		t.Fatalf("expected assertion to fail, but got nil error")
	} else {
		t.Logf("expected assertion to fail, got error: %v", err)
	}
}

func TestEvaluateStatusAndHeader(t *testing.T) {
	resp := &HTTPResponse{
		Status: 200,
		Headers: http.Header{
			"Content-Type": []string{"application/json; charset=utf-8"},
			"X-Request-Id": []string{"abc123"},
			"X-Debug-Flag": []string{"on"},
		},
		Body:     []byte(`{"ok":true}`),
		Duration: 50 * time.Millisecond,
	}

	assertText := `
	status == 200
	header "Content-Type" contains "application/json"
	header "X-Request-Id" exists
	header "X-Missing" not_exists
	`

	if err := Evaluate(assertText, resp); err != nil {
		t.Fatalf("expected status/header assertions to pass, got: %v", err)
	}
}

func TestEvaluateBodyAndDuration(t *testing.T) {
	resp := &HTTPResponse{
		Status:   200,
		Headers:  http.Header{},
		Body:     []byte("Fast response body"),
		Duration: 80 * time.Millisecond,
	}

	assertText := `
	body contains "Fast response"
	duration_ms < 100
	`

	if err := Evaluate(assertText, resp); err != nil {
		t.Fatalf("expected body/duration assertions to pass, got: %v", err)
	}
}

func TestEvaluateGJSONRegexAndExists(t *testing.T) {
	resp := &HTTPResponse{
		Status:   200,
		Headers:  http.Header{},
		Body:     []byte(`{"token":"Abc123_X","nested":{"value":42}}`),
		Duration: 5 * time.Millisecond,
	}

	assertText := `
	gjson "token" matches /^[A-Za-z0-9_]+$/
	gjson "nested.value" > 0
	gjson "nested.value" exists
	gjson "nested.missing" not_exists
	`

	if err := Evaluate(assertText, resp); err != nil {
		t.Fatalf("expected gjson regex/exists assertions to pass, got: %v", err)
	}
}

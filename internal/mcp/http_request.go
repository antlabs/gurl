package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/antlabs/gurl/internal/config"
	"github.com/antlabs/gurl/internal/parser"
	"github.com/mark3labs/mcp-go/mcp"
)

// HTTPResponse represents a structured HTTP response
type HTTPResponse struct {
	Status     string            `json:"status"`
	StatusCode int               `json:"status_code"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
}

// handleHTTPRequest handles the gurl.http_request tool
func (s *Server) handleHTTPRequest(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	defer func() {
		if r := recover(); r != nil {
			Logger.Printf("Panic in handleHTTPRequest: %v", r)
		}
	}()

	// Get URL
	targetURL := mcp.ParseString(req, "url", "")
	if targetURL == "" {
		Logger.Printf("URL is required")
		return nil, fmt.Errorf("url is required")
	}

	Logger.Printf("Processing request to URL: %s", targetURL)

	// Validate and format URL
	if !strings.HasPrefix(targetURL, "http://") && !strings.HasPrefix(targetURL, "https://") {
		targetURL = "http://" + targetURL
	}

	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		Logger.Printf("Invalid URL: %v", err)
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	// Get method
	method := mcp.ParseString(req, "method", "GET")

	// Get headers
	headers := mcp.ParseStringMap(req, "headers", map[string]any{})

	// Get body
	body := mcp.ParseString(req, "body", "")

	// Create config
	cfg := config.Config{
		Method:  method,
		Headers: []string{}, // We'll add headers individually
		Body:    body,
	}

	// Build request
	httpReq, err := parser.BuildRequest(cfg, parsedURL)
	if err != nil {
		Logger.Printf("Failed to build request: %v", err)
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	// Add headers
	for k, v := range headers {
		if vs, ok := v.(string); ok {
			httpReq.Header.Set(k, vs)
		}
	}

	// Execute request
	Logger.Printf("Executing request to %s", targetURL)
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		Logger.Printf("Request failed: %v", err)
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		Logger.Printf("Failed to read response body: %v", err)
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Convert headers to map[string]string
	responseHeaders := make(map[string]string)
	for k, v := range resp.Header {
		responseHeaders[k] = strings.Join(v, ", ")
	}

	// Create structured response
	httpResponse := HTTPResponse{
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		Headers:    responseHeaders,
		Body:       string(bodyBytes),
	}

	// Convert to JSON
	jsonResponse, err := json.Marshal(httpResponse)
	if err != nil {
		Logger.Printf("Failed to marshal response to JSON: %v", err)
		return nil, fmt.Errorf("failed to marshal response to JSON: %w", err)
	}

	Logger.Printf("Request completed successfully, status code: %d", resp.StatusCode)

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Type: "text",
				Text: string(jsonResponse),
			},
		},
	}, nil
}

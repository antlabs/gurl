package mcp

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/antlabs/gurl/internal/config"
	"github.com/antlabs/gurl/internal/parser"
	"github.com/mark3labs/mcp-go/mcp"
)

// handleHTTPRequest handles the gurl.http_request tool
func (s *Server) handleHTTPRequest(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	Logger.Printf("handleHTTPRequest called")
	
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

	// Create a new HTTP client
	client := &http.Client{}

	// Execute request
	Logger.Printf("Executing request to %s", targetURL)
	resp, err := client.Do(httpReq)
	if err != nil {
		Logger.Printf("Request failed: %v", err)
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Format response
	result := fmt.Sprintf("Status: %s\n", resp.Status)
	result += fmt.Sprintf("StatusCode: %d\n", resp.StatusCode)
	
	// Add headers to result
	result += "Headers:\n"
	for k, v := range resp.Header {
		result += fmt.Sprintf("  %s: %s\n", k, strings.Join(v, ", "))
	}
	
	// TODO: Add response body to result
	
	Logger.Printf("Request completed successfully, status code: %d", resp.StatusCode)

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: result,
			},
		},
	}, nil
}
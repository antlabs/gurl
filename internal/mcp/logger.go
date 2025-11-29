package mcp

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/mark3labs/mcp-go/mcp"
)

// Logger is a simple logger for the MCP server
var Logger *log.Logger

func init() {
	// Initialize with default stderr logger
	Logger = log.New(os.Stderr, "[gurl-mcp] ", log.LstdFlags|log.Lshortfile)
}

// InitLogger initializes the logger with optional debug log file
func InitLogger(debugLogFile string) error {
	var writer io.Writer = os.Stderr // Default to stderr

	if debugLogFile != "" {
		// Create directory if it doesn't exist
		dir := filepath.Dir(debugLogFile)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}

		// Open log file for writing (create if not exists, append if exists)
		file, err := os.OpenFile(debugLogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return err
		}

		// Use MultiWriter to write to both stderr and file
		writer = io.MultiWriter(os.Stderr, file)
		log.Printf("MCP debug logging enabled, writing to: %s", debugLogFile)
	}

	// Re-initialize the logger with new writer
	Logger = log.New(writer, "[gurl-mcp] ", log.LstdFlags|log.Lshortfile)
	return nil
}

// ToolHandler represents a MCP tool handler function
type ToolHandler func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error)

// WithLogging wraps a tool handler with request/response logging middleware
func WithLogging(toolName string, handler ToolHandler) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Log the tool call
		Logger.Printf("%s called", toolName)

		// Log request parameters
		logRequestParameters(req)

		// Call the actual handler
		result, err := handler(ctx, req)

		// Log response if successful
		if err == nil && result != nil {
			logResponse(result)
		} else if err != nil {
			Logger.Printf("Tool %s failed with error: %v", toolName, err)
		}

		return result, err
	}
}

// logRequestParameters logs the request parameters in JSON format
func logRequestParameters(req mcp.CallToolRequest) {
	Logger.Printf("Request parameters:")
	if arguments, ok := req.Params.Arguments.(map[string]any); ok && arguments != nil {
		if jsonBytes, err := json.MarshalIndent(arguments, "", "  "); err == nil {
			Logger.Printf("Parameters JSON: %s", string(jsonBytes))
		} else {
			Logger.Printf("Failed to serialize parameters to JSON: %v", err)
			// Fallback to simple format
			for key, value := range arguments {
				Logger.Printf("  %s: %v", key, value)
			}
		}
	} else {
		Logger.Printf("  No parameters provided")
	}
}

// logResponse logs the response in JSON format
func logResponse(result *mcp.CallToolResult) {
	Logger.Printf("Response:%v", result)
	if jsonBytes, err := json.MarshalIndent(result, "", "  "); err == nil {
		Logger.Printf("Response JSON: %s\n", string(jsonBytes))
		Logger.Printf("Result:%v\n", result)
	} else {
		Logger.Printf("Failed to serialize response to JSON: %v", err)
		// Fallback to simple format
		if len(result.Content) > 0 {
			if textContent, ok := result.Content[0].(*mcp.TextContent); ok {
				Logger.Printf("Response text: %s", textContent.Text)
			}
		}
	}
}

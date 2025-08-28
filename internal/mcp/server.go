package mcp

import (
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Server represents the MCP server for gurl
type Server struct {
	mcpServer *server.MCPServer
}

// NewServer creates a new MCP server instance
func NewServer() *Server {
	return &Server{}
}

// Start starts the MCP server
func (s *Server) Start() error {
	// Create a new MCP server
	mcpServer := server.NewMCPServer("gurl", "0.1.0", 
		server.WithToolCapabilities(true))

	// Add tools
	mcpServer.AddTool(
		mcp.NewTool(
			"gurl.http_request",
			mcp.WithDescription("Execute a single HTTP request and return the response"),
			mcp.WithString("url", mcp.Description("The URL to send the request to (required)"), mcp.Required()),
			mcp.WithString("method", mcp.Description("HTTP method (GET, POST, PUT, DELETE, etc.)"), mcp.DefaultString("GET")),
			mcp.WithObject("headers", mcp.Description("HTTP headers to include in the request (key-value pairs)"), mcp.AdditionalProperties(map[string]any{"type": "string"})),
			mcp.WithString("body", mcp.Description("Request body (for POST, PUT, etc.)")),
		), 
		s.handleHTTPRequest,
	)
	
	mcpServer.AddTool(
		mcp.NewTool(
			"gurl.benchmark",
			mcp.WithDescription("Run an HTTP benchmark test and return performance statistics"),
			mcp.WithString("url", mcp.Description("The URL to test (required if not using curl)")),
			mcp.WithString("curl", mcp.Description("Curl command to parse and use for benchmarking (alternative to url)")),
			mcp.WithNumber("connections", mcp.Description("Number of HTTP connections to keep open"), mcp.DefaultNumber(10)),
			mcp.WithString("duration", mcp.Description("Duration of test (e.g., \"30s\", \"1m\")"), mcp.DefaultString("10s")),
			mcp.WithNumber("threads", mcp.Description("Number of threads to use"), mcp.DefaultNumber(2)),
			mcp.WithNumber("rate", mcp.Description("Work rate (requests/sec) 0=unlimited"), mcp.DefaultNumber(0)),
			mcp.WithString("timeout", mcp.Description("Socket/request timeout (e.g., \"5s\")"), mcp.DefaultString("30s")),
			mcp.WithString("method", mcp.Description("HTTP method (if not using curl)"), mcp.DefaultString("GET")),
			mcp.WithObject("headers", mcp.Description("HTTP headers to add to request (if not using curl)"), mcp.AdditionalProperties(map[string]any{"type": "string"})),
			mcp.WithString("body", mcp.Description("HTTP request body (if not using curl)")),
			mcp.WithString("content_type", mcp.Description("Content-Type header (if not using curl)")),
			mcp.WithBoolean("verbose", mcp.Description("Enable verbose output"), mcp.DefaultBool(false)),
			mcp.WithBoolean("latency", mcp.Description("Print detailed latency statistics"), mcp.DefaultBool(false)),
			mcp.WithBoolean("use_nethttp", mcp.Description("Force use standard library net/http instead of pulse"), mcp.DefaultBool(false)),
		),
		s.handleBenchmark,
	)
	
	mcpServer.AddTool(
		mcp.NewTool(
			"gurl.batch_test",
			mcp.WithDescription("Run a batch of tests from a configuration file or inline test definitions and return a summary report"),
			mcp.WithString("config", mcp.Description("Path to batch test configuration file (YAML/JSON)")),
			mcp.WithArray("tests", mcp.Description("Inline batch test definitions")),
			mcp.WithBoolean("verbose", mcp.Description("Enable verbose output"), mcp.DefaultBool(false)),
			mcp.WithNumber("concurrency", mcp.Description("Maximum concurrent batch tests"), mcp.DefaultNumber(3)),
		),
		s.handleBatchTest,
	)

	// Store the server instance
	s.mcpServer = mcpServer

	fmt.Println("gurl MCP server started. Waiting for client connection...")

	// Start the server with stdio transport
	return server.ServeStdio(mcpServer)
}
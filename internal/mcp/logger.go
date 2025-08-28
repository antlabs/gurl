package mcp

import (
	"log"
	"os"
)

// Logger is a simple logger for the MCP server
var Logger *log.Logger

func init() {
	// Initialize the logger to write to stderr
	Logger = log.New(os.Stderr, "[gurl-mcp] ", log.LstdFlags|log.Lshortfile)
}
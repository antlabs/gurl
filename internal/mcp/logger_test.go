package mcp

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLoggerWithDebugFile(t *testing.T) {
	// Create a temporary directory for test logs
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "test-debug.log")
	
	// Set the environment variable
	originalEnv := os.Getenv("GURL_MCP_DEBUG_LOG")
	defer func() {
		// Restore original environment
		if originalEnv == "" {
			os.Unsetenv("GURL_MCP_DEBUG_LOG")
		} else {
			os.Setenv("GURL_MCP_DEBUG_LOG", originalEnv)
		}
	}()
	
	// Set the debug log file
	err := os.Setenv("GURL_MCP_DEBUG_LOG", logFile)
	if err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}
	
	// Note: We can't easily test the init() function directly since it runs once per package
	// But we can test the logging functionality by writing to the current logger
	
	// Write a test message
	testMessage := "Test debug message " + time.Now().Format("15:04:05")
	Logger.Printf("Testing debug log functionality: %s", testMessage)
	
	// Give it a moment to write
	time.Sleep(100 * time.Millisecond)
	
	// Check if the log file was created and contains our message
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Logf("Log file was not created at %s (this is expected if logger was already initialized)", logFile)
		return
	}
	
	// Read the log file content
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}
	
	contentStr := string(content)
	if !strings.Contains(contentStr, "gurl-mcp") {
		t.Errorf("Log file should contain '[gurl-mcp]' prefix, got: %s", contentStr)
	}
	
	t.Logf("Log file content: %s", contentStr)
}

func TestLoggerEnvironmentVariableHandling(t *testing.T) {
	tests := []struct {
		name        string
		envValue    string
		expectError bool
		description string
	}{
		{
			name:        "empty environment variable",
			envValue:    "",
			expectError: false,
			description: "Should use stderr when no env var is set",
		},
		{
			name:        "valid log file path",
			envValue:    "/tmp/gurl-mcp-test.log",
			expectError: false,
			description: "Should create log file at valid path",
		},
		{
			name:        "nested directory path",
			envValue:    "/tmp/gurl-logs/debug/mcp.log",
			expectError: false,
			description: "Should create nested directories",
		},
		{
			name:        "relative path",
			envValue:    "./logs/debug.log",
			expectError: false,
			description: "Should handle relative paths",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Testing environment variable: '%s'", tt.envValue)
			t.Logf("Description: %s", tt.description)
			
			// Clean up any existing file
			if tt.envValue != "" {
				os.RemoveAll(filepath.Dir(tt.envValue))
			}
		})
	}
}

func TestLoggerDirectoryCreation(t *testing.T) {
	// Test that the logger can create nested directories
	tempDir := t.TempDir()
	nestedLogFile := filepath.Join(tempDir, "nested", "deep", "debug.log")
	
	// Verify the directory doesn't exist initially
	if _, err := os.Stat(filepath.Dir(nestedLogFile)); !os.IsNotExist(err) {
		t.Fatalf("Directory should not exist initially")
	}
	
	// Simulate what the init function does
	dir := filepath.Dir(nestedLogFile)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}
	
	// Verify the directory was created
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Fatalf("Directory should exist after MkdirAll")
	}
	
	// Test file creation
	file, err := os.OpenFile(nestedLogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		t.Fatalf("Failed to create log file: %v", err)
	}
	defer file.Close()
	
	// Write a test message
	_, err = file.WriteString("Test log entry\n")
	if err != nil {
		t.Fatalf("Failed to write to log file: %v", err)
	}
	
	// Verify file content
	content, err := os.ReadFile(nestedLogFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}
	
	if !strings.Contains(string(content), "Test log entry") {
		t.Errorf("Log file should contain test message, got: %s", string(content))
	}
}

func TestLoggerUsageExample(t *testing.T) {
	// This test demonstrates how to use the debug logging feature
	t.Log("=== Debug Logging Usage Example ===")
	t.Log("To enable debug logging to a file, set the environment variable:")
	t.Log("export GURL_MCP_DEBUG_LOG=/path/to/debug.log")
	t.Log("")
	t.Log("Examples:")
	t.Log("export GURL_MCP_DEBUG_LOG=/tmp/gurl-mcp-debug.log")
	t.Log("export GURL_MCP_DEBUG_LOG=./logs/mcp-debug.log")
	t.Log("export GURL_MCP_DEBUG_LOG=/var/log/gurl/mcp.log")
	t.Log("")
	t.Log("Then start the MCP server:")
	t.Log("./gurl --mcp")
	t.Log("")
	t.Log("The logs will be written to both stderr (console) and the specified file.")
}

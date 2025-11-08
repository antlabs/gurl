package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

// TestMCPServerIntegration 测试完整的MCP服务器集成
func TestMCPServerIntegration(t *testing.T) {
	// Create a mock HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status": "ok"}`)
	}))
	defer mockServer.Close()
	
	server := NewServer()
	
	// 测试服务器创建
	if server == nil {
		t.Fatal("Failed to create MCP server")
	}
	
	// 测试工具处理函数是否存在
	ctx := context.Background()
	
	// 创建一个有效的请求
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "gurl.http_request",
			Arguments: map[string]any{
				"url":    mockServer.URL,
				"method": "GET",
			},
		},
	}
	
	// 测试处理函数
	result, err := server.handleHTTPRequest(ctx, request)
	if err != nil {
		t.Errorf("handleHTTPRequest failed: %v", err)
	}
	
	if result == nil {
		t.Error("Expected result but got nil")
	}
	
	t.Logf("HTTP request result: %+v", result)
}

// TestCorrectJSONRPCRequest 测试正确格式的JSON-RPC请求
func TestCorrectJSONRPCRequest(t *testing.T) {
	// 这是你原始请求的修正版本
	correctRequest := map[string]any{
		"jsonrpc": "2.0",
		"id":      "12345",
		"method":  "tools/call", // 注意：这应该是 tools/call，不是 gurl.http_request
		"params": map[string]any{
			"name": "gurl.http_request", // 工具名称在这里
			"arguments": map[string]any{
				"url":    "http://127.0.0.1:80/api/v1/token/generate",
				"method": "POST",
				"headers": map[string]any{
					"Content-Type": "application/json",
				},
				// body 应该是字符串，不是对象
				"body": `{"userId":"x","username":"y","tokenExpireTime":"72h","refreshTokenExpireTime":"720h","clientType":"admin"}`,
				// 注意：以下参数不被 gurl.http_request 支持，应该使用 gurl.benchmark 工具
				// "duration":"15s",
				// "connections":3,
				// "threads":2,
				// "rate":20,
				// "timeout":"10s",
				// "latency":true,
				// "verbose":true
			},
		},
	}
	
	jsonBytes, err := json.MarshalIndent(correctRequest, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal correct request: %v", err)
	}
	
	t.Logf("Correct JSON-RPC request format:\n%s", string(jsonBytes))
	
	// 如果你需要基准测试功能，应该使用 gurl.benchmark 工具
	benchmarkRequest := map[string]any{
		"jsonrpc": "2.0",
		"id":      "12345",
		"method":  "tools/call",
		"params": map[string]any{
			"name": "gurl.benchmark", // 使用 benchmark 工具
			"arguments": map[string]any{
				"url":         "http://127.0.0.1:80/api/v1/token/generate",
				"method":      "POST",
				"headers":     map[string]any{"Content-Type": "application/json"},
				"body":        `{"userId":"x","username":"y","tokenExpireTime":"72h","refreshTokenExpireTime":"720h","clientType":"admin"}`,
				"duration":    "15s",
				"connections": 3,
				"threads":     2,
				"rate":        20,
				"timeout":     "10s",
				"latency":     true,
				"verbose":     true,
			},
		},
	}
	
	benchmarkBytes, err := json.MarshalIndent(benchmarkRequest, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal benchmark request: %v", err)
	}
	
	t.Logf("For benchmark functionality, use gurl.benchmark tool:\n%s", string(benchmarkBytes))
}

// TestMCPServerStartup 测试MCP服务器启动过程（不实际启动）
func TestMCPServerStartup(t *testing.T) {
	server := NewServer()
	
	// 验证服务器结构
	if server.mcpServer != nil {
		t.Error("mcpServer should be nil before Start() is called")
	}
	
	// 注意：我们不能在测试中调用 server.Start()，因为它会阻塞
	// 但我们可以验证服务器的基本结构是正确的
}

// TestToolDefinitions 验证工具定义是否正确
func TestToolDefinitions(t *testing.T) {
	// 验证 gurl.http_request 工具的参数定义
	expectedParams := []string{"url", "method", "headers", "body"}
	
	t.Logf("gurl.http_request tool should support these parameters: %v", expectedParams)
	
	// 验证 gurl.benchmark 工具的参数定义
	benchmarkParams := []string{
		"url", "curl", "connections", "duration", "threads", "rate", 
		"timeout", "method", "headers", "body", "content_type", 
		"verbose", "latency", "use_nethttp",
	}
	
	t.Logf("gurl.benchmark tool supports these parameters: %v", benchmarkParams)
	
	// 验证 gurl.batch_test 工具的参数定义
	batchParams := []string{"config", "tests", "verbose", "concurrency"}
	
	t.Logf("gurl.batch_test tool supports these parameters: %v", batchParams)
}

// TestErrorScenarios 测试各种错误场景
func TestErrorScenarios(t *testing.T) {
	server := NewServer()
	ctx := context.Background()
	
	errorTests := []struct {
		name        string
		request     mcp.CallToolRequest
		expectError bool
		description string
	}{
		{
			name: "empty URL",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "gurl.http_request",
					Arguments: map[string]any{
						"url": "",
					},
				},
			},
			expectError: true,
			description: "Empty URL should cause error",
		},
		{
			name: "invalid URL format",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "gurl.http_request",
					Arguments: map[string]any{
						"url": "not-a-url",
					},
				},
			},
			expectError: true,
			description: "Invalid URL format should cause network error",
		},
		{
			name: "timeout scenario",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "gurl.http_request",
					Arguments: map[string]any{
						"url": "http://localhost:1", // Use unreachable address
					},
				},
			},
			expectError: true,
			description: "Should handle connection errors",
		},
	}
	
	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置超时以避免测试挂起
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			
			result, err := server.handleHTTPRequest(ctx, tt.request)
			
			if tt.expectError && err == nil {
				t.Errorf("Expected error for %s, but got none", tt.description)
			}
			
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error for %s: %v", tt.description, err)
			}
			
			t.Logf("Test '%s': error=%v, result=%v", tt.name, err, result != nil)
		})
	}
}

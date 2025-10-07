package mcp

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestHandleHTTPRequest(t *testing.T) {
	server := NewServer()

	tests := []struct {
		name        string
		request     mcp.CallToolRequest
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid simple GET request",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "gurl.http_request",
					Arguments: map[string]any{
						"url":    "http://httpbin.org/get",
						"method": "GET",
					},
				},
			},
			expectError: false,
		},
		{
			name: "valid POST request with headers and body",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "gurl.http_request",
					Arguments: map[string]any{
						"url":    "http://httpbin.org/post",
						"method": "POST",
						"headers": map[string]any{
							"Content-Type": "application/json",
						},
						"body": `{"userId":"x","username":"y"}`,
					},
				},
			},
			expectError: false,
		},
		{
			name: "missing URL should return error",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "gurl.http_request",
					Arguments: map[string]any{
						"method": "GET",
					},
				},
			},
			expectError: true,
			errorMsg:    "url is required",
		},
		{
			name: "invalid URL should return error",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "gurl.http_request",
					Arguments: map[string]any{
						"url":    "not-a-valid-url",
						"method": "GET",
					},
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := server.handleHTTPRequest(context.Background(), tt.request)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result == nil {
					t.Errorf("Expected result but got nil")
				}
			}
		})
	}
}

func TestMCPServerToolRegistration(t *testing.T) {
	// 测试工具是否正确注册
	server := NewServer()

	// 这个测试验证服务器能够正确创建，但不启动实际的服务
	if server == nil {
		t.Fatal("Failed to create MCP server")
	}
}

// TestJSONRPCRequest 测试你提供的具体JSON-RPC请求
func TestJSONRPCRequest(t *testing.T) {
	// 解析你提供的JSON请求
	jsonRequest := `{
		"jsonrpc":"2.0",
		"id":"12345",
		"method":"gurl.http_request",
		"params":{
			"url":"http://127.0.0.1:80/api/v1/token/generate",
			"method":"POST",
			"headers":{"Content-Type":"application/json"},
			"body":{"userId":"x","username":"y","tokenExpireTime":"72h","refreshTokenExpireTime":"720h","clientType":"admin"},
			"duration":"15s",
			"connections":3,
			"threads":2,
			"rate":20,
			"timeout":"10s",
			"latency":true,
			"verbose":true
		}
	}`

	var request map[string]any
	err := json.Unmarshal([]byte(jsonRequest), &request)
	if err != nil {
		t.Fatalf("Failed to parse JSON request: %v", err)
	}

	t.Logf("Parsed request: %+v", request)

	// 分析请求中的问题
	params, ok := request["params"].(map[string]any)
	if !ok {
		t.Fatal("Invalid params structure")
	}

	// 检查body参数 - 这是一个主要问题
	body, exists := params["body"]
	if exists {
		t.Logf("Body type: %T", body)
		t.Logf("Body value: %+v", body)

		// body应该是字符串，但你传递的是对象
		if _, isString := body.(string); !isString {
			t.Logf("PROBLEM: body should be a string, but got %T", body)

			// 正确的做法是将对象序列化为JSON字符串
			bodyBytes, _ := json.Marshal(body)
			t.Logf("Correct body should be: %s", string(bodyBytes))
		}
	}

	// 检查不支持的参数
	unsupportedParams := []string{"duration", "connections", "threads", "rate", "timeout", "latency", "verbose"}
	for _, param := range unsupportedParams {
		if _, exists := params[param]; exists {
			t.Logf("PROBLEM: Parameter '%s' is not supported by gurl.http_request tool", param)
		}
	}
}

func TestHTTPRequest(t *testing.T) {

	resp := HTTPResponse{
		Status:     "200 OK",
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       "{}",
	}

	jsonResponse, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Failed to marshal response to JSON: %v", err)
	}
	t.Logf("JSON response: %s", string(jsonResponse))
	result := &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: string(jsonResponse),
			},
		},
	}

	all, _ := json.Marshal(result)
	t.Logf("%s\n", string(all))
}

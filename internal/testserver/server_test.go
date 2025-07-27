package testserver

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"
)

// TestServerLifecycle 测试服务器生命周期
func TestServerLifecycle(t *testing.T) {
	server := NewTestServer(8080)
	
	// 启动服务器
	go func() {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			t.Errorf("服务器启动失败: %v", err)
		}
	}()
	
	// 等待服务器启动
	time.Sleep(100 * time.Millisecond)
	
	// 测试健康检查
	resp, err := http.Get(server.GetURL() + "/health")
	if err != nil {
		t.Fatalf("健康检查请求失败: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		t.Errorf("期望状态码 200, 得到 %d", resp.StatusCode)
	}
	
	// 停止服务器
	if err := server.Stop(); err != nil {
		t.Errorf("服务器停止失败: %v", err)
	}
}

// TestHTTPMethods 测试各种HTTP方法
func TestHTTPMethods(t *testing.T) {
	server := NewTestServer(8081)
	
	// 启动服务器
	go func() {
		server.Start()
	}()
	defer server.Stop()
	
	// 等待服务器启动
	time.Sleep(100 * time.Millisecond)
	
	baseURL := server.GetURL()
	
	tests := []struct {
		method   string
		endpoint string
		expected int
	}{
		{"GET", "/api/get", http.StatusOK},
		{"POST", "/api/post", http.StatusOK},
		{"PUT", "/api/put", http.StatusOK},
		{"PATCH", "/api/patch", http.StatusOK},
		{"DELETE", "/api/delete", http.StatusOK},
		{"GET", "/api/echo", http.StatusOK},
		{"POST", "/api/echo", http.StatusOK},
	}
	
	for _, test := range tests {
		t.Run(fmt.Sprintf("%s_%s", test.method, test.endpoint), func(t *testing.T) {
			req, err := http.NewRequest(test.method, baseURL+test.endpoint, nil)
			if err != nil {
				t.Fatalf("创建请求失败: %v", err)
			}
			
			client := &http.Client{Timeout: 5 * time.Second}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("请求失败: %v", err)
			}
			defer resp.Body.Close()
			
			if resp.StatusCode != test.expected {
				t.Errorf("期望状态码 %d, 得到 %d", test.expected, resp.StatusCode)
			}
		})
	}
}

// TestDelayEndpoint 测试延迟端点
func TestDelayEndpoint(t *testing.T) {
	server := NewTestServer(8082)
	
	go func() {
		server.Start()
	}()
	defer server.Stop()
	
	time.Sleep(100 * time.Millisecond)
	
	start := time.Now()
	resp, err := http.Get(server.GetURL() + "/api/delay?ms=200")
	duration := time.Since(start)
	
	if err != nil {
		t.Fatalf("延迟请求失败: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		t.Errorf("期望状态码 200, 得到 %d", resp.StatusCode)
	}
	
	// 检查延迟是否大致正确（允许一些误差）
	if duration < 150*time.Millisecond || duration > 300*time.Millisecond {
		t.Errorf("期望延迟约200ms, 实际延迟 %v", duration)
	}
}

// TestStatusEndpoint 测试状态码端点
func TestStatusEndpoint(t *testing.T) {
	server := NewTestServer(8083)
	
	go func() {
		server.Start()
	}()
	defer server.Stop()
	
	time.Sleep(100 * time.Millisecond)
	
	statusCodes := []int{200, 201, 400, 404, 500}
	
	for _, code := range statusCodes {
		t.Run(fmt.Sprintf("status_%d", code), func(t *testing.T) {
			resp, err := http.Get(fmt.Sprintf("%s/api/status/%d", server.GetURL(), code))
			if err != nil {
				t.Fatalf("状态码请求失败: %v", err)
			}
			defer resp.Body.Close()
			
			if resp.StatusCode != code {
				t.Errorf("期望状态码 %d, 得到 %d", code, resp.StatusCode)
			}
		})
	}
}

// getAvailablePort 获取可用端口
func getAvailablePort() int {
	// 简单的端口分配策略，实际使用中可能需要更复杂的逻辑
	return 8080 + int(time.Now().UnixNano()%1000)
}

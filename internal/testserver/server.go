package testserver

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// TestServer 测试服务器结构体
type TestServer struct {
	server *http.Server
	port   int
}

// RequestInfo 请求信息结构体
type RequestInfo struct {
	Method    string            `json:"method"`
	URL       string            `json:"url"`
	Headers   map[string]string `json:"headers"`
	Body      string            `json:"body,omitempty"`
	Timestamp time.Time         `json:"timestamp"`
}

// ResponseData 响应数据结构体
type ResponseData struct {
	Message string      `json:"message"`
	Request RequestInfo `json:"request"`
	Data    interface{} `json:"data,omitempty"`
}

// NewTestServer 创建新的测试服务器
func NewTestServer(port int) *TestServer {
	return &TestServer{
		port: port,
	}
}

// Start 启动测试服务器
func (ts *TestServer) Start() error {
	mux := http.NewServeMux()
	
	// 注册各种HTTP方法的处理器
	mux.HandleFunc("/api/get", ts.handleGet)
	mux.HandleFunc("/api/post", ts.handlePost)
	mux.HandleFunc("/api/put", ts.handlePut)
	mux.HandleFunc("/api/patch", ts.handlePatch)
	mux.HandleFunc("/api/delete", ts.handleDelete)
	
	// 通用处理器，支持所有方法
	mux.HandleFunc("/api/echo", ts.handleEcho)
	
	// 延迟测试端点
	mux.HandleFunc("/api/delay", ts.handleDelay)
	
	// 状态码测试端点
	mux.HandleFunc("/api/status/", ts.handleStatus)
	
	// 健康检查端点
	mux.HandleFunc("/health", ts.handleHealth)
	
	ts.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", ts.port),
		Handler: mux,
	}
	
	log.Printf("测试服务器启动在端口 %d", ts.port)
	return ts.server.ListenAndServe()
}

// Stop 停止测试服务器
func (ts *TestServer) Stop() error {
	if ts.server != nil {
		return ts.server.Close()
	}
	return nil
}

// GetURL 获取服务器URL
func (ts *TestServer) GetURL() string {
	return fmt.Sprintf("http://localhost:%d", ts.port)
}

// handleGet 处理GET请求
func (ts *TestServer) handleGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	response := ResponseData{
		Message: "GET request successful",
		Request: ts.extractRequestInfo(r),
		Data: map[string]interface{}{
			"query_params": r.URL.Query(),
		},
	}
	
	ts.writeJSONResponse(w, response)
}

// handlePost 处理POST请求
func (ts *TestServer) handlePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}
	
	response := ResponseData{
		Message: "POST request successful",
		Request: ts.extractRequestInfo(r),
		Data: map[string]interface{}{
			"body_length": len(body),
			"content_type": r.Header.Get("Content-Type"),
		},
	}
	
	ts.writeJSONResponse(w, response)
}

// handlePut 处理PUT请求
func (ts *TestServer) handlePut(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}
	
	response := ResponseData{
		Message: "PUT request successful",
		Request: ts.extractRequestInfo(r),
		Data: map[string]interface{}{
			"body_length": len(body),
			"updated": true,
		},
	}
	
	ts.writeJSONResponse(w, response)
}

// handlePatch 处理PATCH请求
func (ts *TestServer) handlePatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}
	
	response := ResponseData{
		Message: "PATCH request successful",
		Request: ts.extractRequestInfo(r),
		Data: map[string]interface{}{
			"body_length": len(body),
			"patched": true,
		},
	}
	
	ts.writeJSONResponse(w, response)
}

// handleDelete 处理DELETE请求
func (ts *TestServer) handleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	response := ResponseData{
		Message: "DELETE request successful",
		Request: ts.extractRequestInfo(r),
		Data: map[string]interface{}{
			"deleted": true,
		},
	}
	
	ts.writeJSONResponse(w, response)
}

// handleEcho 通用回显处理器
func (ts *TestServer) handleEcho(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	
	response := ResponseData{
		Message: fmt.Sprintf("%s request to echo endpoint", r.Method),
		Request: ts.extractRequestInfo(r),
		Data: map[string]interface{}{
			"body": string(body),
			"query_params": r.URL.Query(),
		},
	}
	
	ts.writeJSONResponse(w, response)
}

// handleDelay 处理延迟请求
func (ts *TestServer) handleDelay(w http.ResponseWriter, r *http.Request) {
	delayStr := r.URL.Query().Get("ms")
	if delayStr == "" {
		delayStr = "100" // 默认100ms延迟
	}
	
	delay, err := strconv.Atoi(delayStr)
	if err != nil {
		http.Error(w, "Invalid delay parameter", http.StatusBadRequest)
		return
	}
	
	time.Sleep(time.Duration(delay) * time.Millisecond)
	
	response := ResponseData{
		Message: fmt.Sprintf("Delayed response after %dms", delay),
		Request: ts.extractRequestInfo(r),
		Data: map[string]interface{}{
			"delay_ms": delay,
		},
	}
	
	ts.writeJSONResponse(w, response)
}

// handleStatus 处理状态码测试
func (ts *TestServer) handleStatus(w http.ResponseWriter, r *http.Request) {
	statusStr := r.URL.Path[len("/api/status/"):]
	if statusStr == "" {
		statusStr = "200"
	}
	
	status, err := strconv.Atoi(statusStr)
	if err != nil || status < 100 || status > 599 {
		http.Error(w, "Invalid status code", http.StatusBadRequest)
		return
	}
	
	w.WriteHeader(status)
	
	response := ResponseData{
		Message: fmt.Sprintf("Response with status %d", status),
		Request: ts.extractRequestInfo(r),
		Data: map[string]interface{}{
			"status_code": status,
		},
	}
	
	ts.writeJSONResponse(w, response)
}

// handleHealth 健康检查
func (ts *TestServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status": "ok",
		"timestamp": time.Now(),
		"server": "murl-test-server",
	}
	
	ts.writeJSONResponse(w, response)
}

// extractRequestInfo 提取请求信息
func (ts *TestServer) extractRequestInfo(r *http.Request) RequestInfo {
	headers := make(map[string]string)
	for key, values := range r.Header {
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}
	
	body := ""
	if r.Body != nil {
		bodyBytes, _ := io.ReadAll(r.Body)
		body = string(bodyBytes)
		// 重新设置body以便后续读取
		r.Body = io.NopCloser(strings.NewReader(body))
	}
	
	return RequestInfo{
		Method:    r.Method,
		URL:       r.URL.String(),
		Headers:   headers,
		Body:      body,
		Timestamp: time.Now(),
	}
}

// writeJSONResponse 写入JSON响应
func (ts *TestServer) writeJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

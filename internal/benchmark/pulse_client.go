package benchmark

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/antlabs/httparser"
	"github.com/antlabs/murl/internal/config"
	"github.com/antlabs/murl/internal/stats"
	"github.com/antlabs/pulse"
)

// PulseBenchmark 使用pulse库进行HTTP压测的实现
type PulseBenchmark struct {
	config  config.Config
	request *http.Request
	target  *url.URL
}

// HTTPParseResult 存储HTTP解析结果
type HTTPParseResult struct {
	statusCode    int
	contentLength int64
	headersComplete bool
	messageComplete bool
	body          []byte
}

// HTTPClientHandler 处理HTTP客户端连接的回调
type HTTPClientHandler struct {
	request      *http.Request
	requestCount *int64
	errorCount   *int64
	results      *stats.Results
	wg           *sync.WaitGroup
	startTime    time.Time
	responseData []byte
}

// NewPulseBenchmark 创建新的pulse基准测试实例
func NewPulseBenchmark(cfg config.Config, req *http.Request) *PulseBenchmark {
	return &PulseBenchmark{
		config:  cfg,
		request: req,
		target:  req.URL,
	}
}

// OnOpen 连接建立时的回调
func (h *HTTPClientHandler) OnOpen(c *pulse.Conn) {
	// 构建HTTP请求
	httpReq := h.buildHTTPRequest()
	
	// 记录请求开始时间
	h.startTime = time.Now()
	
	// 发送HTTP请求
	_, err := c.Write([]byte(httpReq))
	if err != nil {
		atomic.AddInt64(h.errorCount, 1)
		h.results.AddError(err)
		if h.wg != nil {
			h.wg.Done()
		}
		return
	}
}

// OnData 接收到数据时的回调
func (h *HTTPClientHandler) OnData(c *pulse.Conn, data []byte) {
	// 累积响应数据
	h.responseData = append(h.responseData, data...)
	
	// 检查是否收到完整的HTTP响应
	if h.isCompleteResponse(h.responseData) {
		duration := time.Since(h.startTime)
		atomic.AddInt64(h.requestCount, 1)
		
		// 解析HTTP响应
		statusCode, contentLength := h.parseHTTPResponse(h.responseData)
		
		// 记录统计数据
		h.results.AddLatency(duration)
		h.results.AddStatusCode(statusCode)
		h.results.AddBytes(contentLength)
		
		// 重置响应数据，准备下一个请求
		h.responseData = nil
		
		if h.wg != nil {
			h.wg.Done()
		}
	}
}

// OnClose 连接关闭时的回调
func (h *HTTPClientHandler) OnClose(c *pulse.Conn, err error) {
	if err != nil {
		atomic.AddInt64(h.errorCount, 1)
		h.results.AddError(err)
	}

	if h.wg != nil {
		h.wg.Done()
	}
}

// buildHTTPRequest 构建HTTP请求字符串
func (h *HTTPClientHandler) buildHTTPRequest() string {
	var builder strings.Builder

	// 请求行
	builder.WriteString(fmt.Sprintf("%s %s HTTP/1.1\r\n",
		h.request.Method, h.request.URL.RequestURI()))

	// Host头部
	builder.WriteString(fmt.Sprintf("Host: %s\r\n", h.request.Host))

	// 其他头部
	for name, values := range h.request.Header {
		for _, value := range values {
			builder.WriteString(fmt.Sprintf("%s: %s\r\n", name, value))
		}
	}

	// 如果有请求体，添加Content-Length
	if h.request.Body != nil && h.request.ContentLength > 0 {
		builder.WriteString(fmt.Sprintf("Content-Length: %d\r\n", h.request.ContentLength))
	}

	// 结束头部
	builder.WriteString("\r\n")

	// 添加请求体（如果有）
	if h.request.Body != nil {
		// 读取请求体内容
		body, err := io.ReadAll(h.request.Body)
		if err == nil && len(body) > 0 {
			builder.Write(body)
		}
	}

	return builder.String()
}

// isCompleteResponse 检查是否收到完整的HTTP响应
func (h *HTTPClientHandler) isCompleteResponse(data []byte) bool {
	response := string(data)
	
	// 检查是否包含HTTP响应状态行
	if !strings.Contains(response, "HTTP/") {
		return false
	}
	
	// 检查是否包含完整的头部（以\r\n\r\n结束）
	headerEnd := strings.Index(response, "\r\n\r\n")
	if headerEnd == -1 {
		return false
	}
	
	// 解析Content-Length
	headerPart := response[:headerEnd]
	contentLengthRegex := strings.Contains(strings.ToLower(headerPart), "content-length:")
	if contentLengthRegex {
		// 简化处理：如果有Content-Length，检查body长度
		lines := strings.Split(headerPart, "\r\n")
		for _, line := range lines {
			if strings.HasPrefix(strings.ToLower(line), "content-length:") {
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					contentLength, err := strconv.Atoi(strings.TrimSpace(parts[1]))
					if err == nil {
						body := response[headerEnd+4:]
						return len(body) >= contentLength
					}
				}
			}
		}
	}
	
	// 如果是chunked编码，检查是否以0\r\n\r\n结束
	if strings.Contains(strings.ToLower(headerPart), "transfer-encoding: chunked") {
		return strings.HasSuffix(response, "0\r\n\r\n")
	}
	
	// 默认认为已完整（简化处理）
	return true
}

// parseHTTPResponse 使用httparser解析HTTP响应，返回状态码和内容长度
func (h *HTTPClientHandler) parseHTTPResponse(data []byte) (int, int64) {
	result := &HTTPParseResult{}
	
	// 创建httparser设置
	setting := httparser.Setting{
		MessageBegin: func(*httparser.Parser, int) {
			// 解析器开始工作
		},
		URL: func(_ *httparser.Parser, buf []byte, _ int) {
			// URL数据（响应包中不需要）
		},
		Status: func(_ *httparser.Parser, buf []byte, _ int) {
			// 解析状态码
			if statusCode, err := strconv.Atoi(string(buf)); err == nil {
				result.statusCode = statusCode
			}
		},
		HeaderField: func(_ *httparser.Parser, buf []byte, _ int) {
			// HTTP header field（可以用于调试）
		},
		HeaderValue: func(_ *httparser.Parser, buf []byte, _ int) {
			// HTTP header value（可以用于调试）
		},
		HeadersComplete: func(_ *httparser.Parser, _ int) {
			// HTTP header解析结束
			result.headersComplete = true
		},
		Body: func(_ *httparser.Parser, buf []byte, _ int) {
			// 累积响应体数据
			result.body = append(result.body, buf...)
		},
		MessageComplete: func(_ *httparser.Parser, _ int) {
			// 消息解析结束
			result.messageComplete = true
		},
	}
	
	// 创建HTTP响应解析器
	p := httparser.New(httparser.RESPONSE)
	
	// 执行解析
	success, err := p.Execute(&setting, data)
	if err != nil {
		// 解析失败时使用原始方法作为后备
		return h.parseHTTPResponseFallback(data)
	}
	
	// 如果解析成功，计算内容长度
	contentLength := int64(len(result.body))
	if success > 0 && result.statusCode > 0 {
		return result.statusCode, contentLength
	}
	
	// 如果httparser解析不完整，使用后备方法
	return h.parseHTTPResponseFallback(data)
}

// parseHTTPResponseFallback 原始HTTP解析方法作为后备
func (h *HTTPClientHandler) parseHTTPResponseFallback(data []byte) (int, int64) {
	response := string(data)
	statusCode := 0
	contentLength := int64(len(data))
	
	// 解析状态码
	lines := strings.Split(response, "\r\n")
	if len(lines) > 0 {
		statusLine := lines[0]
		parts := strings.Split(statusLine, " ")
		if len(parts) >= 2 {
			fmt.Sscanf(parts[1], "%d", &statusCode)
		}
	}
	
	// 计算实际内容长度（响应体部分）
	headerEnd := strings.Index(response, "\r\n\r\n")
	if headerEnd != -1 {
		body := response[headerEnd+4:]
		contentLength = int64(len(body))
	}
	
	return statusCode, contentLength
}

// parseStatusCode 从HTTP响应中解析状态码（保持向后兼容）
func (h *HTTPClientHandler) parseStatusCode(response string) int {
	lines := strings.Split(response, "\r\n")
	if len(lines) > 0 {
		statusLine := lines[0]
		parts := strings.Split(statusLine, " ")
		if len(parts) >= 2 {
			var statusCode int
			fmt.Sscanf(parts[1], "%d", &statusCode)
			return statusCode
		}
	}
	return 0
}

// Run 执行pulse基准测试
func (pb *PulseBenchmark) Run(ctx context.Context) (*stats.Results, error) {
	results := stats.NewResults()

	// 创建测试上下文
	testCtx, cancel := context.WithTimeout(ctx, pb.config.Duration)
	defer cancel()

	var requestCount int64
	var errorCount int64
	var wg sync.WaitGroup

	// 创建pulse客户端事件循环
	loop := pulse.NewClientEventLoop(
		testCtx,
		pulse.WithCallback(&HTTPClientHandler{
			request:      pb.request,
			requestCount: &requestCount,
			errorCount:   &errorCount,
			results:      results,
			wg:           &wg,
		}),
	)

	// 建立连接
	port := pb.target.Port()
	if port == "" {
		if pb.target.Scheme == "https" {
			port = "443"
		} else {
			port = "80"
		}
	}

	address := fmt.Sprintf("%s:%s", pb.target.Hostname(), port)

	// 创建多个连接
	for i := 0; i < pb.config.Connections; i++ {
		conn, err := net.Dial("tcp", address)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to %s: %w", address, err)
		}

		wg.Add(1)
		loop.RegisterConn(conn)
	}

	// 启动事件循环
	go func() {
		loop.Serve()
	}()

	// 等待所有请求完成或超时
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-testCtx.Done():
	}

	// 计算最终结果
	results.TotalRequests = atomic.LoadInt64(&requestCount)
	results.TotalErrors = atomic.LoadInt64(&errorCount)
	results.Duration = pb.config.Duration

	return results, nil
}

// ShouldUsePulse 判断是否应该使用pulse库
func ShouldUsePulse(req *http.Request) bool {
	// 只对HTTP（非HTTPS）请求使用pulse
	return req.URL.Scheme == "http"
}

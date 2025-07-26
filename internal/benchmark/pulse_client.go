package benchmark

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
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
	statusCode      int
	contentLength   int64
	headersComplete bool
	messageComplete bool
	body            []byte
}

// HTTPClientHandler 处理HTTP客户端连接的回调
type HTTPClientHandler struct {
	request      *http.Request
	requestCount *int64
	errorCount   *int64
	results      *stats.Results
	wg           *sync.WaitGroup
	startTime    time.Time
	maxBodySize  int64

	// 流式解析相关
	parser      *httparser.Parser
	parseResult *HTTPParseResult
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
	// 初始化解析器
	h.parseResult = &HTTPParseResult{}
	h.parser = httparser.New(httparser.RESPONSE)
	h.parser.SetUserData(h.parseResult)

	// 记录请求开始时间
	h.startTime = time.Now()

	// 构建HTTP请求
	httpReq := h.buildHTTPRequest()

	// 发送HTTP请求
	_, err := c.Write(httpReq)
	if err != nil {
		atomic.AddInt64(h.errorCount, 1)
		h.results.AddError(err)
		h.wg.Done()
		c.Close()
		return
	}
}

// OnData 接收到数据时的回调
func (h *HTTPClientHandler) OnData(c *pulse.Conn, data []byte) {
	// 流式解析HTTP响应
	n, err := h.parser.Execute(&httpParserSetting, data)
	if err != nil || n < 0 {
		atomic.AddInt64(h.errorCount, 1)
		h.results.AddError(fmt.Errorf("HTTP parse error: %v", err))
		h.wg.Done()
		c.Close()
		return
	}

	// 检查是否收到完整的HTTP响应
	if h.parseResult.messageComplete {
		duration := time.Since(h.startTime)
		atomic.AddInt64(h.requestCount, 1)

		// 记录统计数据
		h.results.AddLatency(duration)
		h.results.AddStatusCode(h.parseResult.statusCode)
		h.results.AddBytes(h.parseResult.contentLength)

		// 重置解析器状态，准备下次请求
		h.parseResult = &HTTPParseResult{}
		h.parser.SetUserData(h.parseResult)

		h.wg.Done()
	}
}

// OnClose 连接关闭时的回调
func (h *HTTPClientHandler) OnClose(c *pulse.Conn, err error) {
	if err != nil {
		atomic.AddInt64(h.errorCount, 1)
		h.results.AddError(err)
	}

	h.wg.Done()
}

// buildHTTPRequest 构建HTTP请求字符串
func (h *HTTPClientHandler) buildHTTPRequest() []byte {
	b, err := httputil.DumpRequest(h.request, true)
	if err != nil {
		return nil
	}
	return b
}

// httpParserSetting 全局HTTP解析器设置，避免每次创建
var httpParserSetting = httparser.Setting{
	MessageBegin: func(p *httparser.Parser, _ int) {
		// 解析器开始工作，重置结果
		if result := p.GetUserData(); result != nil {
			if r, ok := result.(*HTTPParseResult); ok {
				r.statusCode = 0
				r.contentLength = 0
				r.headersComplete = false
				r.messageComplete = false
				r.body = r.body[:0] // 重用slice，避免重新分配
			}
		}
	},
	URL: func(_ *httparser.Parser, buf []byte, _ int) {
		// URL数据（响应包中不需要）
	},
	Status: func(p *httparser.Parser, buf []byte, _ int) {
		// 解析状态码
		if result := p.GetUserData(); result != nil {
			if r, ok := result.(*HTTPParseResult); ok {
				if statusCode, err := strconv.Atoi(string(buf)); err == nil {
					r.statusCode = statusCode
				}
			}
		}
	},
	HeaderField: func(p *httparser.Parser, buf []byte, _ int) {
		// HTTP header field

		if string(buf) == "Content-Length" {
			if result := p.GetUserData(); result != nil {
				if r, ok := result.(*HTTPParseResult); ok {
					if l, err := strconv.ParseInt(string(buf), 10, 64); err == nil {
						r.contentLength = l
					}
				}
			}
		}

	},
	HeaderValue: func(p *httparser.Parser, buf []byte, _ int) {
		// HTTP header value
	},
	HeadersComplete: func(p *httparser.Parser, _ int) {
		// HTTP header解析结束
		if result := p.GetUserData(); result != nil {
			if r, ok := result.(*HTTPParseResult); ok {
				r.headersComplete = true
			}
		}
	},
	Body: func(p *httparser.Parser, buf []byte, _ int) {
		// 累积响应体数据并检查大小限制
		if result := p.GetUserData(); result != nil {
			if r, ok := result.(*HTTPParseResult); ok {
				r.body = append(r.body, buf...)
				// 检查body大小限制（例如1MB）
				if int64(len(r.body)) > 1<<20 {
					// 超过限制，记录错误并关闭连接
					// 这里可以通过设置标志位在OnData中处理
				}
			}
		}
	},
	MessageComplete: func(p *httparser.Parser, _ int) {
		// 消息解析结束
		if result := p.GetUserData(); result != nil {
			if r, ok := result.(*HTTPParseResult); ok {
				r.messageComplete = true
				// 如果之前没有设置contentLength，则使用body长度
				if r.contentLength == 0 {
					r.contentLength = int64(len(r.body))
				}
			}
		}
	},
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
			maxBodySize:  1 << 20, // 1MB限制
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

	address := net.JoinHostPort(pb.target.Hostname(), port)

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

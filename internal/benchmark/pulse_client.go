package benchmark

import (
	"bytes"
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
	"github.com/antlabs/pulse/core"
)

var bytesContentLength = []byte("Content-Length")

// PulseBenchmark 使用pulse库进行HTTP压测的实现
type PulseBenchmark struct {
	config  config.Config
	request *http.Request
	target  *url.URL
}

// HTTPParseResult 存储HTTP解析结果
type HTTPParseResult struct {
	statusCode       int
	contentLength    int64
	headersComplete  bool
	messageComplete  bool
	hasContentLength bool
	// body             []byte
}

func (h *HTTPParseResult) Reset() {
	h.statusCode = 0
	h.contentLength = 0
	h.headersComplete = false
	h.messageComplete = false
	h.hasContentLength = false
}

// ConnSession 每个连接的会话状态
type ConnSession struct {
	startTime    time.Time
	parser       *httparser.Parser
	parseResult  *HTTPParseResult
	request      *http.Request
	requestCount *int64
	errorCount   *int64
	results      *stats.Results
	wg           *sync.WaitGroup
	maxBodySize  int64
}

// HTTPClientHandler 处理HTTP客户端连接的回调
type HTTPClientHandler struct {
	request      *http.Request
	requestCount *int64
	errorCount   *int64
	results      *stats.Results
	wg           *sync.WaitGroup
	maxBodySize  int64
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
	session := &ConnSession{
		startTime:    time.Now(),
		parseResult:  &HTTPParseResult{},
		request:      h.request,
		requestCount: h.requestCount,
		errorCount:   h.errorCount,
		results:      h.results,
		wg:           h.wg,
		maxBodySize:  h.maxBodySize,
	}

	session.parser = httparser.New(httparser.RESPONSE)
	session.parser.SetUserData(session.parseResult)

	c.SetSession(session)

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
	session, ok := c.GetSession().(*ConnSession)
	if !ok {
		return
	}

	// 流式解析HTTP响应
	n, err := session.parser.Execute(&httpParserSetting, data)
	if err != nil || n < 0 {
		atomic.AddInt64(h.errorCount, 1)
		session.results.AddError(fmt.Errorf("HTTP parse error: %v", err))
		c.Close()
		return
	}

	// 检查是否收到完整的HTTP响应
	if session.parseResult.messageComplete {
		duration := time.Since(session.startTime)
		atomic.AddInt64(h.requestCount, 1)

		// 记录统计数据
		session.results.AddLatency(duration)
		session.results.AddStatusCode(session.parseResult.statusCode)
		session.results.AddBytes(session.parseResult.contentLength)

		// 重置解析器状态，准备下次请求
		session.parseResult.Reset()
		session.parser.SetUserData(session.parseResult)
		session.startTime = time.Now()

		// 立即发送下一个请求（持续压测）
		httpReq := h.buildHTTPRequest()
		_, writeErr := c.Write(httpReq)
		if writeErr != nil {
			atomic.AddInt64(h.errorCount, 1)
			session.results.AddError(writeErr)
			c.Close()
			return
		}
	}
}

// OnClose 连接关闭时的回调
func (h *HTTPClientHandler) OnClose(c *pulse.Conn, err error) {
	session, ok := c.GetSession().(*ConnSession)
	if !ok {
		return
	}

	if err != nil {
		atomic.AddInt64(h.errorCount, 1)
		session.results.AddError(err)
	}

	// 只有在连接关闭时才调用 wg.Done()
	session.wg.Done()
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
				r.contentLength = 0
				// r.body = r.body[:0] // 重用slice，避免重新分配
			}
		}
	},
	URL: func(_ *httparser.Parser, buf []byte, _ int) {
	},
	Status: func(p *httparser.Parser, buf []byte, _ int) {
	},
	HeaderField: func(p *httparser.Parser, buf []byte, _ int) {
		// HTTP header field
		if bytes.Equal(buf, bytesContentLength) {
			if result := p.GetUserData(); result != nil {
				if r, ok := result.(*HTTPParseResult); ok {
					r.hasContentLength = true
				}
			}
		}
	},
	HeaderValue: func(p *httparser.Parser, buf []byte, _ int) {
		// HTTP header value
		// 获取当前header field
		if result := p.GetUserData(); result != nil {
			if r, ok := result.(*HTTPParseResult); ok {
				if r.hasContentLength {
					if contentLength, err := strconv.Atoi(string(buf)); err == nil {
						r.contentLength = int64(contentLength)
					}
				}
			}
		}
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
				// r.body = append(r.body, buf...)
				r.contentLength += int64(len(buf))
				// 检查body大小限制（例如1MB）
				if r.contentLength > 1<<20 {
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
				r.statusCode = int(p.StatusCode)
				r.messageComplete = true
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
	// var wg sync.WaitGroup

	// 创建pulse客户端事件循环
	loop := pulse.NewClientEventLoop(
		testCtx,
		pulse.WithTaskType(pulse.TaskTypeInEventLoop), // 在事件循环中处理任务
		pulse.WithTriggerType(core.TriggerTypeLevel),
		pulse.WithCallback(&HTTPClientHandler{
			request:      pb.request,
			requestCount: &requestCount,
			errorCount:   &errorCount,
			results:      results,
			// wg:           &wg,
			maxBodySize: 1 << 20, // 1MB限制
		}),
	)

	// 启动事件循环
	go func() {
		loop.Serve()
	}()
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
	fmt.Printf("Creating connections...%s\n", address)
	for i := 0; i < pb.config.Connections; i++ {
		conn, err := net.Dial("tcp", address)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to %s: %w", address, err)
		}

		err = loop.RegisterConn(conn)
		if err != nil {
			return nil, fmt.Errorf("failed to register connection: %w", err)
		}
		// wg.Add(1)
	}

	// 等待测试时间结束
	<-testCtx.Done()

	fmt.Println("testCtx.Done()")
	// 等待所有连接关闭（通过context取消触发）
	// wg.Wait()
	fmt.Println("wg.Wait()")

	// 计算最终结果
	results.TotalRequests = atomic.LoadInt64(&requestCount)
	results.TotalErrors = atomic.LoadInt64(&errorCount)
	results.Duration = pb.config.Duration

	return results, nil
}

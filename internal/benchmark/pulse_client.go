package benchmark

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/antlabs/gurl/internal/asserts"
	"github.com/antlabs/gurl/internal/config"
	"github.com/antlabs/gurl/internal/stats"
	"github.com/antlabs/httparser"
	"github.com/antlabs/pulse"
	"github.com/antlabs/pulse/core"
	"go.uber.org/ratelimit"
)

var bytesContentLength = []byte("Content-Length")

// PulseBenchmark 使用 pulse 库进行单请求 HTTP 压测的实现
type PulseBenchmark struct {
	config  config.Config
	request *http.Request
	target  *url.URL
}

// PulseBenchmarkMulti 使用 pulse 库进行多请求 HTTP 压测的实现
// 多请求列表通过 RequestPool 管理，按照配置的负载策略分发
type PulseBenchmarkMulti struct {
	config      config.Config
	requestPool *RequestPool
	target      *url.URL
}

// HTTPParseResult 存储HTTP解析结果
type HTTPParseResult struct {
	statusCode       int
	contentLength    int64
	headersComplete  bool
	messageComplete  bool
	hasContentLength bool
	// 为 asserts 收集的数据
	enableAsserts bool
	maxBodySize   int64
	headers       http.Header
	body          []byte
	currentHeader string
}

func (h *HTTPParseResult) Reset() {
	h.statusCode = 0
	h.contentLength = 0
	h.headersComplete = false
	h.messageComplete = false
	h.hasContentLength = false
	if h.enableAsserts {
		for k := range h.headers {
			delete(h.headers, k)
		}
		if h.body != nil {
			h.body = h.body[:0]
		}
		h.currentHeader = ""
	}
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
	maxBodySize  int64
}

// HTTPClientHandler 处理HTTP客户端连接的回调
type HTTPClientHandler struct {
	request      *http.Request // 单请求模式使用
	requestPool  *RequestPool  // 多请求模式使用
	requestCount *int64
	errorCount   *int64
	results      *stats.Results
	maxBodySize  int64
	rateLimiter  ratelimit.Limiter
	asserts      string
	maxRequests  int64
	cancel       context.CancelFunc
}

// NewPulseBenchmark 创建新的pulse基准测试实例
func NewPulseBenchmark(cfg config.Config, req *http.Request) *PulseBenchmark {
	return &PulseBenchmark{
		config:  cfg,
		request: req,
		target:  req.URL,
	}
}

// NewPulseBenchmarkWithMultipleRequests 创建支持多请求的 pulse 基准测试实例
// 目前假定所有请求指向同一主机，仅使用第一个请求的 URL 建立底层 TCP 连接
func NewPulseBenchmarkWithMultipleRequests(cfg config.Config, requests []*http.Request) *PulseBenchmarkMulti {
	if len(requests) == 0 {
		return nil
	}

	pool := NewRequestPool(requests, cfg.LoadStrategy)
	return &PulseBenchmarkMulti{
		config:      cfg,
		requestPool: pool,
		target:      requests[0].URL,
	}
}

// OnOpen 连接建立时的回调
func (h *HTTPClientHandler) OnOpen(c *pulse.Conn) {
	session := &ConnSession{
		startTime: time.Now(),
		parseResult: &HTTPParseResult{
			enableAsserts: h.asserts != "",
			maxBodySize:   h.maxBodySize,
		},
		request:      h.request,
		requestCount: h.requestCount,
		errorCount:   h.errorCount,
		results:      h.results,
		maxBodySize:  h.maxBodySize,
	}

	session.parser = httparser.New(httparser.RESPONSE)
	session.parser.SetUserData(session.parseResult)

	c.SetSession(session)

	// 在发送第一个请求前，根据 MaxRequests 使用 CAS 占用一个请求名额
	if h.maxRequests > 0 {
		for {
			cur := atomic.LoadInt64(h.requestCount)
			if cur >= h.maxRequests {
				// 已达到上限，取消测试并关闭连接
				if h.cancel != nil {
					h.cancel()
				}
				c.Close()
				return
			}
			if atomic.CompareAndSwapInt64(h.requestCount, cur, cur+1) {
				// 如果这是最后一个名额，占用后立即取消上下文
				if cur+1 >= h.maxRequests && h.cancel != nil {
					h.cancel()
				}
				break
			}
		}
	}

	// 如果配置了频率限制，在发送第一个请求前获取令牌
	if h.rateLimiter != nil {
		h.rateLimiter.Take()
	}

	// 构建HTTP请求（单请求或多请求）
	httpReq := h.buildHTTPRequest()

	// 发送HTTP请求
	written, err := c.Write(httpReq)
	if err != nil {
		atomic.AddInt64(h.errorCount, 1)
		h.results.AddError(err)
		c.Close()
		return
	}

	h.results.AddWriteBytes(int64(written))
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

		// 如果没有配置 maxRequests（=0），在每次完成响应时递增请求计数
		if h.maxRequests == 0 {
			atomic.AddInt64(h.requestCount, 1)
		}

		// 记录统计数据
		session.results.AddLatency(duration)
		session.results.AddStatusCode(session.parseResult.statusCode)
		session.results.AddBytes(session.parseResult.contentLength)

		// 如果配置了断言，则执行断言
		if h.asserts != "" && session.parseResult.enableAsserts {
			assertResp := &asserts.HTTPResponse{
				Status:   session.parseResult.statusCode,
				Headers:  session.parseResult.headers,
				Body:     session.parseResult.body,
				Duration: duration,
			}

			if errAssert := asserts.Evaluate(h.asserts, assertResp); errAssert != nil {
				atomic.AddInt64(h.errorCount, 1)
				session.results.AddError(errAssert)
			}
		}

		// 重置解析器状态，准备下次请求
		session.parseResult.Reset()
		session.parser.SetUserData(session.parseResult)
		session.startTime = time.Now()

		// 在发送下一个请求前，根据 MaxRequests 使用 CAS 占用一个请求名额
		if h.maxRequests > 0 {
			for {
				cur := atomic.LoadInt64(h.requestCount)
				if cur >= h.maxRequests {
					// 已达到上限，取消测试并关闭连接
					if h.cancel != nil {
						h.cancel()
					}
					c.Close()
					return
				}
				if atomic.CompareAndSwapInt64(h.requestCount, cur, cur+1) {
					// 如果这是最后一个名额，占用后立即取消上下文
					if cur+1 >= h.maxRequests && h.cancel != nil {
						h.cancel()
					}
					break
				}
			}
		}

		// 在发送下一个请求前应用频率限制（如果配置了Rate）
		if h.rateLimiter != nil {
			h.rateLimiter.Take()
		}

		// 立即发送下一个请求（持续压测）
		httpReq := h.buildHTTPRequest()
		written, writeErr := c.Write(httpReq)
		if writeErr != nil {
			atomic.AddInt64(h.errorCount, 1)
			session.results.AddError(writeErr)
			c.Close()
			return
		}

		// 记录写入字节数（请求体+头部），优先使用 ContentLength
		session.results.AddWriteBytes(int64(written))
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
	// session.wg.Done()
}

// buildHTTPRequest 构建HTTP请求字符串
func (h *HTTPClientHandler) buildHTTPRequest() []byte {
	var req *http.Request
	if h.requestPool != nil {
		// 多请求模式：从请求池中获取下一个请求
		var size int
		req, size = h.requestPool.GetRequest()
		_ = size // 目前仅用于选择请求，写入字节数使用实际 written 统计
	} else {
		// 单请求模式
		req = h.request
	}

	b, err := httputil.DumpRequest(req, true)
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
				if r.enableAsserts {
					if r.headers == nil {
						r.headers = make(http.Header)
					} else {
						for k := range r.headers {
							delete(r.headers, k)
						}
					}
					if r.body != nil {
						r.body = r.body[:0]
					}
					r.currentHeader = ""
				}
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
		if result := p.GetUserData(); result != nil {
			if r, ok := result.(*HTTPParseResult); ok && r.enableAsserts {
				r.currentHeader = string(buf)
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
				if r.enableAsserts && r.currentHeader != "" {
					if r.headers == nil {
						r.headers = make(http.Header)
					}
					r.headers.Add(r.currentHeader, string(buf))
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
				r.contentLength += int64(len(buf))
				if r.enableAsserts && r.maxBodySize > 0 {
					remaining := int(r.maxBodySize) - len(r.body)
					if remaining > 0 {
						if remaining < len(buf) {
							r.body = append(r.body, buf[:remaining]...)
						} else {
							r.body = append(r.body, buf...)
						}
					}
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

	startTime := time.Now()

	// 创建测试上下文
	testCtx, cancel := context.WithTimeout(ctx, pb.config.Duration)
	defer cancel()

	var requestCount int64
	var errorCount int64

	// 创建 Uber 限流器（所有连接共享，与 NetHTTPBenchmark 行为一致）
	var limiter ratelimit.Limiter
	if pb.config.Rate > 0 {
		limiter = ratelimit.New(pb.config.Rate)
	}

	// 创建 Live UI（如果启用）
	var liveUI *LiveUI
	var uiErr error
	if pb.config.LiveUI {
		liveUI, uiErr = NewLiveUIWithTheme(pb.config.Duration, pb.config.UITheme)
		if uiErr != nil {
			// 如果 UI 初始化失败，继续运行但不显示 UI
			pb.config.LiveUI = false
		} else {
			defer liveUI.Close()
		}
	}

	// 创建 pulse 客户端事件循环
	loop := pulse.NewClientEventLoop(
		testCtx,
		pulse.WithTaskType(pulse.TaskTypeInEventLoop), // 在事件循环中处理任务
		pulse.WithTriggerType(core.TriggerTypeLevel),
		pulse.WithLogLevel(slog.LevelError), // 只显示错误日志，避免INFO日志干扰UI显示
		pulse.WithCallback(&HTTPClientHandler{
			request:      pb.request,
			requestPool:  nil,
			requestCount: &requestCount,
			errorCount:   &errorCount,
			results:      results,
			maxBodySize:  1 << 20, // 1MB限制
			rateLimiter:  limiter,
			asserts:      pb.config.Asserts,
			maxRequests:  pb.config.Requests,
			cancel:       cancel,
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

	// 启动采样 goroutine，每秒记录请求数和更新 UI（在连接建立之前启动）
	samplingDone := StartSampling(testCtx, cancel, &requestCount, &errorCount, results, liveUI, nil, startTime)

	// 创建多个连接（不输出日志，避免破坏 UI）
	for i := 0; i < pb.config.Connections; i++ {
		conn, err := net.Dial("tcp", address)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to %s: %w", address, err)
		}

		err = loop.RegisterConn(conn)
		if err != nil {
			return nil, fmt.Errorf("failed to register connection: %w", err)
		}
	}

	// 等待测试时间结束
	<-testCtx.Done()

	// 等待采样完成
	<-samplingDone

	// 计算最终结果
	results.TotalRequests = atomic.LoadInt64(&requestCount)
	results.TotalErrors = atomic.LoadInt64(&errorCount)
	// 使用实际运行时间，而不是配置的时间（支持提前中断）
	results.Duration = time.Since(startTime)

	return results, nil
}

// Run 执行 pulse 多请求基准测试
func (pb *PulseBenchmarkMulti) Run(ctx context.Context) (*stats.Results, error) {
	results := stats.NewResults()

	startTime := time.Now()

	// 创建测试上下文
	testCtx, cancel := context.WithTimeout(ctx, pb.config.Duration)
	defer cancel()

	var requestCount int64
	var errorCount int64

	// 创建 Uber 限流器（所有连接共享，与 NetHTTPBenchmark 行为一致）
	var limiter ratelimit.Limiter
	if pb.config.Rate > 0 {
		limiter = ratelimit.New(pb.config.Rate)
	}

	// 创建 Live UI（如果启用）
	var liveUI *LiveUI
	var uiErr error
	if pb.config.LiveUI {
		liveUI, uiErr = NewLiveUIWithTheme(pb.config.Duration, pb.config.UITheme)
		if uiErr != nil {
			// 如果 UI 初始化失败，继续运行但不显示 UI
			pb.config.LiveUI = false
		} else {
			defer liveUI.Close()
		}
	}

	// 创建 pulse 客户端事件循环
	loop := pulse.NewClientEventLoop(
		testCtx,
		pulse.WithTaskType(pulse.TaskTypeInEventLoop), // 在事件循环中处理任务
		pulse.WithTriggerType(core.TriggerTypeLevel),
		pulse.WithLogLevel(slog.LevelError), // 只显示错误日志，避免 INFO 日志干扰 UI 显示
		pulse.WithCallback(&HTTPClientHandler{
			request:      nil,
			requestPool:  pb.requestPool,
			requestCount: &requestCount,
			errorCount:   &errorCount,
			results:      results,
			maxBodySize:  1 << 20, // 1MB 限制
			rateLimiter:  limiter,
			asserts:      pb.config.Asserts,
			maxRequests:  pb.config.Requests,
			cancel:       cancel,
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

	// 启动采样 goroutine，每秒记录请求数和更新 UI（在连接建立之前启动）
	samplingDone := StartSampling(testCtx, cancel, &requestCount, &errorCount, results, liveUI, nil, startTime)

	// 创建多个连接（不输出日志，避免破坏 UI）
	for i := 0; i < pb.config.Connections; i++ {
		conn, err := net.Dial("tcp", address)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to %s: %w", address, err)
		}

		if err := loop.RegisterConn(conn); err != nil {
			return nil, fmt.Errorf("failed to register connection: %w", err)
		}
	}

	// 等待测试时间结束
	<-testCtx.Done()

	// 等待采样完成
	<-samplingDone

	// 计算最终结果
	results.TotalRequests = atomic.LoadInt64(&requestCount)
	results.TotalErrors = atomic.LoadInt64(&errorCount)
	// 使用实际运行时间，而不是配置的时间（支持提前中断）
	results.Duration = time.Since(startTime)

	return results, nil
}

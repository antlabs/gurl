package benchmark

import (
	"context"
	"io"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/antlabs/gurl/internal/asserts"
	"github.com/antlabs/gurl/internal/config"
	"github.com/antlabs/gurl/internal/stats"
	"go.uber.org/ratelimit"
)

// Runner 定义基准测试运行器接口
type Runner interface {
	Run(ctx context.Context) (*stats.Results, error)
}

// NetHTTPBenchmark represents a net/http based benchmark instance
type NetHTTPBenchmark struct {
	config      config.Config
	request     *http.Request
	bodyContent string       // 存储请求体内容，用于每次创建新的 Body
	requestPool *RequestPool // 多请求池
	client      *http.Client
	rateLimiter ratelimit.Limiter // Uber 限流器
}

// NewNetHTTPBenchmark creates a new net/http benchmark instance
func NewNetHTTPBenchmark(cfg config.Config, req *http.Request) *NetHTTPBenchmark {
	// 创建HTTP客户端
	client := &http.Client{
		Timeout: cfg.Timeout,
		Transport: &http.Transport{
			MaxIdleConns:        cfg.Connections,
			MaxIdleConnsPerHost: cfg.Connections,
			IdleConnTimeout:     30 * time.Second,
		},
	}

	// 创建 Uber 限流器（所有连接共享）
	var limiter ratelimit.Limiter
	if cfg.Rate > 0 {
		limiter = ratelimit.New(cfg.Rate) // 每秒请求数
	}

	// 读取并保存 body 内容
	var bodyContent string
	if req.Body != nil {
		bodyBytes, _ := io.ReadAll(req.Body)
		req.Body.Close()
		bodyContent = string(bodyBytes)
		// 重新设置 Body 以便第一次使用
		req.Body = io.NopCloser(strings.NewReader(bodyContent))
	}

	return &NetHTTPBenchmark{
		config:      cfg,
		request:     req,
		bodyContent: bodyContent,
		requestPool: nil, // 单请求模式
		client:      client,
		rateLimiter: limiter,
	}
}

// NewNetHTTPBenchmarkWithMultipleRequests creates a new net/http benchmark with multiple requests
func NewNetHTTPBenchmarkWithMultipleRequests(cfg config.Config, requests []*http.Request) *NetHTTPBenchmark {
	// 创建HTTP客户端
	client := &http.Client{
		Timeout: cfg.Timeout,
		Transport: &http.Transport{
			MaxIdleConns:        cfg.Connections,
			MaxIdleConnsPerHost: cfg.Connections,
			IdleConnTimeout:     30 * time.Second,
		},
	}

	// 创建 Uber 限流器（所有连接共享）
	var limiter ratelimit.Limiter
	if cfg.Rate > 0 {
		limiter = ratelimit.New(cfg.Rate) // 每秒请求数
	}

	// 创建请求池
	requestPool := NewRequestPool(requests, cfg.LoadStrategy)

	return &NetHTTPBenchmark{
		config:      cfg,
		request:     nil, // 多请求模式不使用单个请求
		requestPool: requestPool,
		client:      client,
		rateLimiter: limiter,
	}
}

// Run executes the net/http benchmark
func (b *NetHTTPBenchmark) Run(ctx context.Context) (*stats.Results, error) {
	results := stats.NewResults()

	// 创建上下文，在指定时间后取消；如果配置了最大请求数，也会在达到后取消
	testCtx, cancel := context.WithTimeout(ctx, b.config.Duration)
	defer cancel()

	var wg sync.WaitGroup
	var requestCount int64
	var errorCount int64

	// 初始化 Live UI（如果启用）
	var liveUI *LiveUI
	var uiErr error
	if b.config.LiveUI {
		liveUI, uiErr = NewLiveUIWithTheme(b.config.Duration, b.config.UITheme)
		if uiErr != nil {
			// 如果 UI 初始化失败，继续运行但不显示 UI
			b.config.LiveUI = false
		} else {
			defer liveUI.Close()
		}
	}

	// 记录开始时间
	startTime := time.Now()

	// 启动采样 goroutine，每秒记录请求数
	samplingDone := StartSampling(testCtx, cancel, &requestCount, &errorCount, results, liveUI, b.requestPool, startTime)

	// 启动工作线程
	for i := 0; i < b.config.Threads; i++ {
		wg.Add(1)
		go func(threadID int) {
			defer wg.Done()
			b.runWorker(testCtx, cancel, threadID, &requestCount, &errorCount, results)
		}(i)
	}

	// 等待所有工作线程完成
	wg.Wait()

	// 等待采样完成
	<-samplingDone

	// 计算最终结果
	results.TotalRequests = atomic.LoadInt64(&requestCount)
	results.TotalErrors = atomic.LoadInt64(&errorCount)
	// 使用实际运行时间，而不是配置的时间（支持提前中断）
	results.Duration = time.Since(startTime)

	return results, nil
}

// runWorker runs a single worker thread
func (b *NetHTTPBenchmark) runWorker(ctx context.Context, cancel context.CancelFunc, threadID int, requestCount, errorCount *int64, results *stats.Results) {
	connectionsPerThread := b.config.Connections / b.config.Threads
	if threadID < b.config.Connections%b.config.Threads {
		connectionsPerThread++
	}

	var wg sync.WaitGroup

	// 为每个连接启动一个goroutine
	for i := 0; i < connectionsPerThread; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			b.runConnection(ctx, cancel, requestCount, errorCount, results)
		}()
	}
	wg.Wait()
}

// runConnection handles a single connection's requests
func (b *NetHTTPBenchmark) runConnection(ctx context.Context, cancel context.CancelFunc, requestCount, errorCount *int64, results *stats.Results) {
	for {
		// 先检查 context 是否已取消
		select {
		case <-ctx.Done():
			return
		default:
		}

		// 如果配置了最大请求数（Requests > 0），使用 CAS 控制总请求次数
		if b.config.Requests > 0 {
			for {
				cur := atomic.LoadInt64(requestCount)
				if cur >= b.config.Requests {
					// 已达到上限，取消测试并退出
					if cancel != nil {
						cancel()
					}
					return
				}
				// 尝试占用一个请求名额
				if atomic.CompareAndSwapInt64(requestCount, cur, cur+1) {
					// 如果这是最后一个名额，占用后立即取消上下文，通知其它 goroutine 退出
					if cur+1 >= b.config.Requests && cancel != nil {
						cancel()
					}
					break
				}
			}
		}

		// Uber 限流器：自动等待并获取令牌
		// Take() 会阻塞，所以在调用前先检查 context
		if b.rateLimiter != nil {
			// 在 goroutine 中调用 Take()，这样可以响应 context 取消
			done := make(chan struct{})
			go func() {
				b.rateLimiter.Take()
				close(done)
			}()

			select {
			case <-done:
				// 获得令牌，继续执行
			case <-ctx.Done():
				// context 取消，立即返回
				return
			}
		}

		// 获取要执行的请求
		var req *http.Request
		if b.requestPool != nil {
			// 多请求模式：从请求池获取
			req = b.requestPool.GetRequest()
		} else {
			// 单请求模式
			req = b.request
		}

		// 克隆请求并创建新的 Body（避免数据竞争）
		clonedReq := req.Clone(ctx)
		if b.bodyContent != "" {
			clonedReq.Body = io.NopCloser(strings.NewReader(b.bodyContent))
		}

		// 执行请求
		start := time.Now()
		resp, err := b.client.Do(clonedReq)
		duration := time.Since(start)

		// 如果没有配置 Requests（=0），使用简单的每请求计数
		if b.config.Requests == 0 {
			atomic.AddInt64(requestCount, 1)
		}

		var bytesRead int64
		var statusCode int

		if err != nil {
			atomic.AddInt64(errorCount, 1)
			results.AddError(err)
		} else {
			// 根据是否配置了断言决定是否需要读取响应体
			if b.config.Asserts != "" {
				// 需要做断言：读取完整响应体，保留 headers 和 duration
				bodyBytes, _ := io.ReadAll(resp.Body)
				_ = resp.Body.Close()
				bytesRead = int64(len(bodyBytes))
				statusCode = resp.StatusCode

				results.AddLatency(duration)
				results.AddStatusCode(statusCode)
				results.AddBytes(bytesRead)

				// 构造断言所需的 HTTPResponse 并执行断言
				assertResp := &asserts.HTTPResponse{
					Status:   statusCode,
					Headers:  resp.Header,
					Body:     bodyBytes,
					Duration: duration,
				}

				if errAssert := asserts.Evaluate(b.config.Asserts, assertResp); errAssert != nil {
					atomic.AddInt64(errorCount, 1)
					results.AddError(errAssert)
				}
			} else {
				// 未配置断言：保持原有高性能行为，仅统计字节数
				bytesRead, _ = io.Copy(io.Discard, resp.Body)
				_ = resp.Body.Close()
				statusCode = resp.StatusCode

				results.AddLatency(duration)
				results.AddStatusCode(statusCode)
				results.AddBytes(bytesRead)
			}
		}

		// 如果是多请求模式，记录每个 URL 的统计
		if b.requestPool != nil {
			results.AddLatencyWithURL(req.URL.String(), duration, statusCode, bytesRead, err)
		}
	}
}

package benchmark

import (
	"context"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/antlabs/murl/internal/config"
	"github.com/antlabs/murl/internal/stats"
)

// Runner 定义基准测试运行器接口
type Runner interface {
	Run(ctx context.Context) (*stats.Results, error)
}

// NetHTTPBenchmark represents a net/http based benchmark instance
type NetHTTPBenchmark struct {
	config  config.Config
	request *http.Request
	client  *http.Client
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

	return &NetHTTPBenchmark{
		config:  cfg,
		request: req,
		client:  client,
	}
}



// Run executes the net/http benchmark
func (b *NetHTTPBenchmark) Run(ctx context.Context) (*stats.Results, error) {
	results := stats.NewResults()
	
	// 创建上下文，在指定时间后取消
	testCtx, cancel := context.WithTimeout(ctx, b.config.Duration)
	defer cancel()

	var wg sync.WaitGroup
	var requestCount int64
	var errorCount int64

	// 启动工作线程
	for i := 0; i < b.config.Threads; i++ {
		wg.Add(1)
		go func(threadID int) {
			defer wg.Done()
			b.runWorker(testCtx, threadID, &requestCount, &errorCount, results)
		}(i)
	}

	// 等待所有工作线程完成
	wg.Wait()

	// 计算最终结果
	results.TotalRequests = atomic.LoadInt64(&requestCount)
	results.TotalErrors = atomic.LoadInt64(&errorCount)
	results.Duration = b.config.Duration
	
	return results, nil
}

// runWorker runs a single worker thread
func (b *NetHTTPBenchmark) runWorker(ctx context.Context, threadID int, requestCount, errorCount *int64, results *stats.Results) {
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
			b.runConnection(ctx, requestCount, errorCount, results)
		}()
	}
	
	wg.Wait()
}

// runConnection handles a single connection's requests
func (b *NetHTTPBenchmark) runConnection(ctx context.Context, requestCount, errorCount *int64, results *stats.Results) {
	rateLimiter := b.createRateLimiter()
	
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// 速率限制
			if rateLimiter != nil {
				select {
				case <-rateLimiter:
				case <-ctx.Done():
					return
				}
			}
			
			// 执行请求
			start := time.Now()
			resp, err := b.client.Do(b.request.Clone(ctx))
			duration := time.Since(start)
			
			atomic.AddInt64(requestCount, 1)
			
			if err != nil {
				atomic.AddInt64(errorCount, 1)
				results.AddError(err)
			} else {
				// 读取并丢弃响应体数据，计算字节数
				bytesRead, _ := io.Copy(io.Discard, resp.Body)
				resp.Body.Close()
				
				results.AddLatency(duration)
				results.AddStatusCode(resp.StatusCode)
				results.AddBytes(bytesRead)
			}
		}
	}
}

// createRateLimiter creates a rate limiter if rate limiting is enabled
func (b *NetHTTPBenchmark) createRateLimiter() <-chan time.Time {
	if b.config.Rate <= 0 {
		return nil
	}
	
	// 计算每个连接的速率
	ratePerConnection := b.config.Rate / b.config.Connections
	if ratePerConnection <= 0 {
		ratePerConnection = 1
	}
	
	interval := time.Second / time.Duration(ratePerConnection)
	return time.Tick(interval)
}

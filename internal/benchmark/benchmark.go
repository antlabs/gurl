package benchmark

import (
	"context"
	"net/http"

	"github.com/antlabs/gurl/internal/config"
	"github.com/antlabs/gurl/internal/stats"
)

// Benchmark 统一的基准测试接口
type Benchmark struct {
	runner Runner
}

// New 创建新的基准测试实例，根据URL自动选择实现
func New(cfg config.Config, req *http.Request) *Benchmark {
	var runner Runner

	// 如果用户强制使用标准库，则使用NetHTTP实现
	if cfg.UseNetHTTP {
		runner = NewNetHTTPBenchmark(cfg, req)
	} else if ShouldUsePulse(req) {
		// 根据URL和请求特性选择合适的实现
		runner = NewPulseBenchmark(cfg, req)
	} else {
		runner = NewNetHTTPBenchmark(cfg, req)
	}

	return &Benchmark{
		runner: runner,
	}
}

// NewWithMultipleRequests 创建支持多请求的基准测试实例
func NewWithMultipleRequests(cfg config.Config, requests []*http.Request) *Benchmark {
	// 多请求模式目前只支持 NetHTTP
	runner := NewNetHTTPBenchmarkWithMultipleRequests(cfg, requests)

	return &Benchmark{
		runner: runner,
	}
}

// Run 执行基准测试
func (b *Benchmark) Run(ctx context.Context) (*stats.Results, error) {
	return b.runner.Run(ctx)
}

// ShouldUsePulse 判断是否应该使用pulse库
// 这个函数已经在pulse_client.go中定义，这里只是为了保持接口一致性
func ShouldUsePulse(req *http.Request) bool {
	// 对于HTTP/HTTPS请求，如果需要高性能或特殊处理，使用pulse
	// 目前的逻辑在pulse_client.go中已经实现
	return shouldUsePulseImpl(req)
}

// shouldUsePulseImpl 实际的判断逻辑（避免循环导入）
func shouldUsePulseImpl(req *http.Request) bool {
	// 这里可以根据具体需求调整判断逻辑
	// 例如：特定的URL模式、请求方法、头部等

	// 重新启用pulse实现进行测试
	scheme := req.URL.Scheme
	return scheme == "http" || scheme == "https"
}

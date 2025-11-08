package benchmark

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/antlabs/gurl/internal/config"
)

// BenchmarkNetHTTPClient 测试NetHTTP客户端性能
func BenchmarkNetHTTPClient(b *testing.B) {
	// 创建mock服务器
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer mockServer.Close()

	// 创建配置
	cfg := config.Config{
		Connections: 1,
		Duration:    100 * time.Millisecond,
		Threads:     1,
		UseNetHTTP:  true,
	}

	// 创建请求
	req, _ := http.NewRequest("GET", mockServer.URL, nil)

	// 创建benchmark实例
	bench := NewNetHTTPBenchmark(cfg, req)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := bench.Run(context.Background())
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkHTTPRequestParsing 测试HTTP请求构建性能
func BenchmarkHTTPRequestParsing(b *testing.B) {
	req, _ := http.NewRequest("POST", "http://example.com/api", nil)
	req.Header.Set("Content-Type", "application/json")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = req.Clone(context.Background())
	}
}

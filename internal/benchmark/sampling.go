package benchmark

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/antlabs/gurl/internal/stats"
)

// StartSampling 启动采样 goroutine，每秒记录请求数并更新 UI
// 返回一个 channel，当采样完成时会关闭
func StartSampling(
	ctx context.Context,
	cancel context.CancelFunc,
	requestCount *int64,
	errorCount *int64,
	results *stats.Results,
	liveUI *LiveUI,
	requestPool *RequestPool,
	startTime time.Time,
) chan struct{} {
	samplingDone := make(chan struct{})

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		lastCount := int64(0)
		for {
			select {
			case <-ticker.C:
				currentCount := atomic.LoadInt64(requestCount)
				reqThisSecond := currentCount - lastCount
				results.AddReqPerSecond(reqThisSecond)
				lastCount = currentCount

				// 更新 Live UI
				if liveUI != nil {
					avgLatency := results.GetAverageLatency()
					minLatency := results.GetMinLatency()
					maxLatency := results.GetMaxLatency()
					statusCodes := results.GetStatusCodes()
					latencyPercentiles := results.GetLatencyPercentiles()
					errors := atomic.LoadInt64(errorCount)
					liveUI.Update(currentCount, reqThisSecond, statusCodes, avgLatency, minLatency, maxLatency, latencyPercentiles, errors)

					// 如果是多端点模式，更新每个端点的统计
					if requestPool != nil {
						endpointStats := results.GetEndpointStats()
						elapsed := time.Since(startTime)
						if elapsed == 0 {
							elapsed = time.Second
						}

						for url, stats := range endpointStats {
							reqPerSec := float64(stats.Requests) / elapsed.Seconds()
							avgLat := stats.GetAverageLatency()
							liveUI.UpdateEndpointStats(url, stats.Requests, reqPerSec, avgLat, stats.MinLatency, stats.MaxLatency, stats.Errors)
						}
					}

					liveUI.Render()
				}
			case <-ctx.Done():
				// 记录最后一个不完整的时间段
				currentCount := atomic.LoadInt64(requestCount)
				if currentCount > lastCount {
					results.AddReqPerSecond(currentCount - lastCount)
				}
				close(samplingDone)
				return
			case <-func() <-chan struct{} {
				if liveUI != nil {
					return liveUI.StopChan()
				}
				return nil
			}():
				// 用户按下退出键，提前取消测试
				cancel()
				close(samplingDone)
				return
			}
		}
	}()

	return samplingDone
}

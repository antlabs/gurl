package benchmark

import (
	"fmt"
	"sort"
	"time"

	"github.com/antlabs/gurl/internal/config"
	"github.com/antlabs/gurl/internal/stats"
)

// PrintResults prints the benchmark results in wrk-like format
func PrintResults(results *stats.Results, cfg config.Config) {
	fmt.Printf("  Thread Stats   Avg      Stdev     Max   +/- Stdev\n")

	// 计算延迟统计
	latencies := results.GetLatencies()
	if len(latencies) > 0 {
		sort.Slice(latencies, func(i, j int) bool {
			return latencies[i] < latencies[j]
		})

		avg := results.GetAverageLatency()
		stdev := results.GetLatencyStdDev()
		max := latencies[len(latencies)-1]

		fmt.Printf("    Latency   %8s %8s %8s %8.2f%%\n",
			formatDuration(avg),
			formatDuration(stdev),
			formatDuration(max),
			calculateStdDevPercentage(latencies, avg, stdev))
	}

	// 计算 Req/Sec 统计
	qps := float64(results.TotalRequests) / results.Duration.Seconds()
	avg, stdev, max, percentage := results.GetReqPerSecStats()
	if avg > 0 {
		fmt.Printf("    Req/Sec   %8.2f %8.2f %8.2f %8.2f%%\n", avg, stdev, max, percentage)
	} else {
		// 如果没有采样数据，使用总体 QPS
		fmt.Printf("    Req/Sec   %8.2f %8s %8s %8s\n", qps, "N/A", "N/A", "N/A")
	}

	// 打印延迟分布
	if cfg.PrintLatency && len(latencies) > 0 {
		fmt.Printf("  Latency Distribution\n")
		percentiles := []float64{50, 75, 90, 99}
		for _, p := range percentiles {
			idx := int(float64(len(latencies)-1) * p / 100.0)
			fmt.Printf("     %2.0f%%   %s\n", p, formatDuration(latencies[idx]))
		}
	}

	// 打印总体统计
	fmt.Printf("  %d requests in %s, %s read\n",
		results.TotalRequests,
		results.Duration,
		formatBytes(results.GetTotalBytes()))

	if results.TotalErrors > 0 {
		fmt.Printf("  Socket errors: connect %d, read %d, write %d, timeout %d\n",
			results.GetConnectErrors(),
			results.GetReadErrors(),
			results.GetWriteErrors(),
			results.GetTimeoutErrors())
	}

	// 打印状态码分布
	statusCodes := results.GetStatusCodes()
	if len(statusCodes) > 0 {
		fmt.Printf("  Status code distribution:\n")
		for code, count := range statusCodes {
			percentage := float64(count) / float64(results.TotalRequests) * 100
			fmt.Printf("    [%d] %d responses (%.1f%%)\n", code, count, percentage)
		}
	}

	fmt.Printf("Requests/sec: %8.2f\n", qps)
	fmt.Printf("Transfer/sec: %8s\n", formatBytes(int64(float64(results.GetTotalBytes())/results.Duration.Seconds())))

	// 打印每个端点的统计（如果有多个端点）
	endpointStats := results.GetEndpointStats()
	if len(endpointStats) > 1 {
		fmt.Printf("\n=== Per-Endpoint Statistics ===\n")

		// 按 URL 排序以保持输出一致性
		urls := make([]string, 0, len(endpointStats))
		for url := range endpointStats {
			urls = append(urls, url)
		}
		sort.Strings(urls)

		for _, url := range urls {
			stats := endpointStats[url]
			printEndpointStats(stats, results.Duration)
		}
	}
}

// printEndpointStats prints statistics for a single endpoint
func printEndpointStats(stats *stats.EndpointStats, duration time.Duration) {
	fmt.Printf("\n[%s]\n", stats.URL)

	// 基本统计
	fmt.Printf("  Requests:     %d\n", stats.Requests)
	if stats.Errors > 0 {
		errorRate := float64(stats.Errors) / float64(stats.Requests) * 100
		fmt.Printf("  Errors:       %d (%.1f%%)\n", stats.Errors, errorRate)
	}

	// TPS
	tps := float64(stats.Requests) / duration.Seconds()
	fmt.Printf("  Requests/sec: %.2f\n", tps)

	// 延迟统计
	if len(stats.Latencies) > 0 {
		avgLatency := stats.GetAverageLatency()
		fmt.Printf("  Latency:      avg=%s, min=%s, max=%s\n",
			formatDuration(avgLatency),
			formatDuration(stats.MinLatency),
			formatDuration(stats.MaxLatency))
	}

	// 状态码分布
	if len(stats.StatusCodes) > 0 {
		fmt.Printf("  Status codes: ")
		first := true
		for code, count := range stats.StatusCodes {
			if !first {
				fmt.Printf(", ")
			}
			percentage := float64(count) / float64(stats.Requests) * 100
			fmt.Printf("[%d] %d (%.1f%%)", code, count, percentage)
			first = false
		}
		fmt.Printf("\n")
	}

	// 数据传输
	if stats.TotalBytes > 0 {
		fmt.Printf("  Data:         %s total, %s/sec\n",
			formatBytes(stats.TotalBytes),
			formatBytes(int64(float64(stats.TotalBytes)/duration.Seconds())))
	}
}

// formatDuration formats a duration for display
func formatDuration(d time.Duration) string {
	if d < time.Microsecond {
		return fmt.Sprintf("%.2fns", float64(d.Nanoseconds()))
	} else if d < time.Millisecond {
		return fmt.Sprintf("%.2fus", float64(d.Nanoseconds())/1000.0)
	} else if d < time.Second {
		return fmt.Sprintf("%.2fms", float64(d.Nanoseconds())/1000000.0)
	} else {
		return fmt.Sprintf("%.2fs", d.Seconds())
	}
}

// formatBytes formats bytes for display
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%dB", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// calculateStdDevPercentage calculates the percentage of values within one standard deviation
func calculateStdDevPercentage(latencies []time.Duration, avg, stdev time.Duration) float64 {
	if len(latencies) == 0 || stdev == 0 {
		return 0
	}

	lower := avg - stdev
	upper := avg + stdev
	count := 0

	for _, lat := range latencies {
		if lat >= lower && lat <= upper {
			count++
		}
	}

	return float64(count) / float64(len(latencies)) * 100.0
}

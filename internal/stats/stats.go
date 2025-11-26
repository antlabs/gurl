package stats

import (
	"math"
	"sort"
	"sync"
	"time"
)

// EndpointStats holds statistics for a single endpoint
type EndpointStats struct {
	URL         string
	Requests    int64
	Errors      int64
	Latencies   []time.Duration
	StatusCodes map[int]int64
	ReadBytes   int64
	WriteBytes  int64
	MinLatency  time.Duration
	MaxLatency  time.Duration
}

// Results holds benchmark results
type Results struct {
	mu              sync.RWMutex
	latencies       []time.Duration
	statusCodes     map[int]int64
	errors          []error
	totalReadBytes  int64
	totalWriteBytes int64
	reqPerSecond    []int64       // 每秒的请求数统计
	minLatency      time.Duration // 最快响应时间
	maxLatency      time.Duration // 最慢响应时间

	// 按 URL 分组的统计
	endpointStats map[string]*EndpointStats

	TotalRequests int64
	TotalErrors   int64
	Duration      time.Duration
}

// NewResults creates a new Results instance
func NewResults() *Results {
	return &Results{
		latencies:     make([]time.Duration, 0),
		statusCodes:   make(map[int]int64),
		errors:        make([]error, 0),
		reqPerSecond:  make([]int64, 0),
		endpointStats: make(map[string]*EndpointStats),
	}
}

// AddLatency adds a latency measurement
func (r *Results) AddLatency(latency time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.latencies = append(r.latencies, latency)

	// 更新最小和最大延迟
	if r.minLatency == 0 || latency < r.minLatency {
		r.minLatency = latency
	}
	if latency > r.maxLatency {
		r.maxLatency = latency
	}
}

// AddLatencyWithURL adds a latency measurement for a specific URL
func (r *Results) AddLatencyWithURL(url string, latency time.Duration, statusCode int, bytes int64, writeBytes int64, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 全局统计
	r.latencies = append(r.latencies, latency)
	if r.minLatency == 0 || latency < r.minLatency {
		r.minLatency = latency
	}
	if latency > r.maxLatency {
		r.maxLatency = latency
	}

	// 按 URL 统计
	if r.endpointStats[url] == nil {
		r.endpointStats[url] = &EndpointStats{
			URL:         url,
			Latencies:   make([]time.Duration, 0),
			StatusCodes: make(map[int]int64),
		}
	}

	stats := r.endpointStats[url]
	stats.Requests++
	stats.Latencies = append(stats.Latencies, latency)
	stats.ReadBytes += bytes
	stats.WriteBytes += writeBytes

	if err != nil {
		stats.Errors++
	} else {
		stats.StatusCodes[statusCode]++
	}

	if stats.MinLatency == 0 || latency < stats.MinLatency {
		stats.MinLatency = latency
	}
	if latency > stats.MaxLatency {
		stats.MaxLatency = latency
	}
}

// AddStatusCode adds a status code count
func (r *Results) AddStatusCode(code int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.statusCodes[code]++
}

// AddError adds an error
func (r *Results) AddError(err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.errors = append(r.errors, err)
}

// AddBytes adds to the total bytes transferred
func (r *Results) AddBytes(bytes int64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.totalReadBytes += bytes
}

// AddWriteBytes adds to the total bytes written (request bodies)
func (r *Results) AddWriteBytes(bytes int64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.totalWriteBytes += bytes
}

// GetLatencies returns a copy of all latencies
func (r *Results) GetLatencies() []time.Duration {
	r.mu.RLock()
	defer r.mu.RUnlock()

	latencies := make([]time.Duration, len(r.latencies))
	copy(latencies, r.latencies)
	return latencies
}

// GetStatusCodes returns a copy of status code counts
func (r *Results) GetStatusCodes() map[int]int64 {
	r.mu.RLock()
	defer r.mu.RUnlock()

	codes := make(map[int]int64)
	for k, v := range r.statusCodes {
		codes[k] = v
	}
	return codes
}

// GetErrors returns a copy of all errors
func (r *Results) GetErrors() []error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	errors := make([]error, len(r.errors))
	copy(errors, r.errors)
	return errors
}

// GetTotalBytes returns the total bytes transferred
func (r *Results) GetTotalBytes() int64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.totalReadBytes
}

// GetTotalWriteBytes returns the total bytes written (request bodies)
func (r *Results) GetTotalWriteBytes() int64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.totalWriteBytes
}

// GetAverageLatency calculates the average latency
func (r *Results) GetAverageLatency() time.Duration {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.latencies) == 0 {
		return 0
	}

	var total time.Duration
	for _, lat := range r.latencies {
		total += lat
	}

	return total / time.Duration(len(r.latencies))
}

// GetLatencyStdDev calculates the standard deviation of latencies
func (r *Results) GetLatencyStdDev() time.Duration {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.latencies) <= 1 {
		return 0
	}

	avg := r.getAverageLatencyUnsafe()
	var sumSquares float64

	for _, lat := range r.latencies {
		diff := float64(lat - avg)
		sumSquares += diff * diff
	}

	variance := sumSquares / float64(len(r.latencies)-1)
	return time.Duration(math.Sqrt(variance))
}

// getAverageLatencyUnsafe calculates average latency without locking (internal use)
func (r *Results) getAverageLatencyUnsafe() time.Duration {
	if len(r.latencies) == 0 {
		return 0
	}

	var total time.Duration
	for _, lat := range r.latencies {
		total += lat
	}

	return total / time.Duration(len(r.latencies))
}

// GetConnectErrors returns the number of connection errors
func (r *Results) GetConnectErrors() int64 {
	// 这里可以根据具体的错误类型进行分类
	// 暂时返回0，实际实现需要解析错误类型
	return 0
}

// GetReadErrors returns the number of read errors
func (r *Results) GetReadErrors() int64 {
	return 0
}

// GetWriteErrors returns the number of write errors
func (r *Results) GetWriteErrors() int64 {
	return 0
}

// GetTimeoutErrors returns the number of timeout errors
func (r *Results) GetTimeoutErrors() int64 {
	return 0
}

// AddReqPerSecond adds a request per second sample
func (r *Results) AddReqPerSecond(count int64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.reqPerSecond = append(r.reqPerSecond, count)
}

// GetReqPerSecond returns a copy of all req/sec samples
func (r *Results) GetReqPerSecond() []int64 {
	r.mu.RLock()
	defer r.mu.RUnlock()

	samples := make([]int64, len(r.reqPerSecond))
	copy(samples, r.reqPerSecond)
	return samples
}

// GetReqPerSecStats calculates statistics for req/sec
func (r *Results) GetReqPerSecStats() (avg, stdev, max float64, percentage float64) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.reqPerSecond) == 0 {
		return 0, 0, 0, 0
	}

	// 计算平均值
	var sum int64
	max = 0
	for _, v := range r.reqPerSecond {
		sum += v
		if float64(v) > max {
			max = float64(v)
		}
	}
	avg = float64(sum) / float64(len(r.reqPerSecond))

	// 计算标准差
	if len(r.reqPerSecond) <= 1 {
		return avg, 0, max, 0
	}

	var sumSquares float64
	for _, v := range r.reqPerSecond {
		diff := float64(v) - avg
		sumSquares += diff * diff
	}
	variance := sumSquares / float64(len(r.reqPerSecond)-1)
	stdev = math.Sqrt(variance)

	// 计算在一个标准差范围内的百分比
	lower := avg - stdev
	upper := avg + stdev
	count := 0
	for _, v := range r.reqPerSecond {
		if float64(v) >= lower && float64(v) <= upper {
			count++
		}
	}
	percentage = float64(count) / float64(len(r.reqPerSecond)) * 100.0

	return avg, stdev, max, percentage
}

// GetMinLatency returns the minimum latency
func (r *Results) GetMinLatency() time.Duration {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.minLatency
}

// GetMaxLatency returns the maximum latency
func (r *Results) GetMaxLatency() time.Duration {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.maxLatency
}

// GetLatencyPercentiles returns latency percentiles (p50, p75, p90, p95, p99)
// 使用采样以避免对大数据集排序
func (r *Results) GetLatencyPercentiles() map[float64]time.Duration {
	r.mu.RLock()
	defer r.mu.RUnlock()

	percentiles := map[float64]time.Duration{}

	if len(r.latencies) == 0 {
		return percentiles
	}

	// 采样：最多使用 10000 个样本
	sampleSize := len(r.latencies)
	if sampleSize > 10000 {
		sampleSize = 10000
	}

	// 使用最近的样本
	startIdx := len(r.latencies) - sampleSize
	sample := make([]time.Duration, sampleSize)
	copy(sample, r.latencies[startIdx:])

	// 使用标准库排序（快速排序，O(n log n)）
	sort.Slice(sample, func(i, j int) bool {
		return sample[i] < sample[j]
	})

	// 计算百分位
	ps := []float64{50, 75, 90, 95, 99}
	for _, p := range ps {
		idx := int(float64(len(sample)-1) * p / 100.0)
		percentiles[p] = sample[idx]
	}

	return percentiles
}

// GetEndpointStats returns statistics for all endpoints
func (r *Results) GetEndpointStats() map[string]*EndpointStats {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// 返回副本以避免并发问题
	result := make(map[string]*EndpointStats)
	for url, stats := range r.endpointStats {
		result[url] = &EndpointStats{
			URL:         stats.URL,
			Requests:    stats.Requests,
			Errors:      stats.Errors,
			Latencies:   append([]time.Duration{}, stats.Latencies...),
			StatusCodes: make(map[int]int64),
			ReadBytes:   stats.ReadBytes,
			WriteBytes:  stats.WriteBytes,
			MinLatency:  stats.MinLatency,
			MaxLatency:  stats.MaxLatency,
		}
		for code, count := range stats.StatusCodes {
			result[url].StatusCodes[code] = count
		}
	}
	return result
}

// GetAverageLatencyForEndpoint returns the average latency for a specific endpoint
func (stats *EndpointStats) GetAverageLatency() time.Duration {
	if len(stats.Latencies) == 0 {
		return 0
	}

	var total time.Duration
	for _, latency := range stats.Latencies {
		total += latency
	}
	return total / time.Duration(len(stats.Latencies))
}

package stats

import (
	"math"
	"sync"
	"time"
)

// Results holds benchmark results
type Results struct {
	mu            sync.RWMutex
	latencies     []time.Duration
	statusCodes   map[int]int64
	errors        []error
	totalBytes    int64
	
	TotalRequests int64
	TotalErrors   int64
	Duration      time.Duration
}

// NewResults creates a new Results instance
func NewResults() *Results {
	return &Results{
		latencies:   make([]time.Duration, 0),
		statusCodes: make(map[int]int64),
		errors:      make([]error, 0),
	}
}

// AddLatency adds a latency measurement
func (r *Results) AddLatency(latency time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.latencies = append(r.latencies, latency)
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
	r.totalBytes += bytes
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
	return r.totalBytes
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

package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/antlabs/gurl/internal/benchmark"
	"github.com/antlabs/gurl/internal/config"
	"github.com/antlabs/gurl/internal/parser"
	"github.com/antlabs/gurl/internal/stats"
)

// BenchmarkRequest represents a benchmark request
type BenchmarkRequest struct {
	URL         string                 `json:"url"`
	Curl        string                 `json:"curl,omitempty"`
	Connections int                    `json:"connections,omitempty"`
	Duration    string                 `json:"duration,omitempty"`
	Threads     int                    `json:"threads,omitempty"`
	Rate        int                    `json:"rate,omitempty"`
	Timeout     string                 `json:"timeout,omitempty"`
	Method      string                 `json:"method,omitempty"`
	Headers     map[string]string      `json:"headers,omitempty"`
	Body        string                 `json:"body,omitempty"`
	ContentType string                 `json:"content_type,omitempty"`
	UseNetHTTP  bool                   `json:"use_nethttp,omitempty"`
	Extra       map[string]interface{} `json:"extra,omitempty"`
}

// BatchRequest represents a batch test request
type BatchRequest struct {
	Tests []BatchTestRequest `json:"tests"`
}

// BatchTestRequest represents a single test in a batch
type BatchTestRequest struct {
	Name        string            `json:"name"`
	Curl        string            `json:"curl,omitempty"`
	URL         string            `json:"url,omitempty"`
	Connections int               `json:"connections,omitempty"`
	Duration    string            `json:"duration,omitempty"`
	Threads     int               `json:"threads,omitempty"`
	Rate        int               `json:"rate,omitempty"`
	Timeout     string            `json:"timeout,omitempty"`
	Method      string            `json:"method,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
	Body        string            `json:"body,omitempty"`
	ContentType string            `json:"content_type,omitempty"`
	UseNetHTTP  bool              `json:"use_nethttp,omitempty"`
}

// BenchmarkResponse represents a benchmark response
type BenchmarkResponse struct {
	TaskID  string                 `json:"task_id"`
	Status  string                 `json:"status"`
	Message string                 `json:"message,omitempty"`
	Config  map[string]interface{} `json:"config,omitempty"`
}

// TaskStatusResponse represents a task status response
type TaskStatusResponse struct {
	ID          string                 `json:"id"`
	Status      string                 `json:"status"`
	CreatedAt   time.Time              `json:"created_at"`
	StartedAt   *time.Time             `json:"started_at,omitempty"`
	CompletedAt *time.Time             `json:"completed_at,omitempty"`
	Error       string                 `json:"error,omitempty"`
	Config      map[string]interface{} `json:"config,omitempty"`
}

// TaskResultsResponse represents task results response
type TaskResultsResponse struct {
	ID          string                `json:"id"`
	Status      string                `json:"status"`
	CreatedAt   time.Time             `json:"created_at"`
	StartedAt   *time.Time            `json:"started_at,omitempty"`
	CompletedAt *time.Time            `json:"completed_at,omitempty"`
	Error       string                `json:"error,omitempty"`
	Results     *BenchmarkResultsJSON `json:"results,omitempty"`
}

// BenchmarkResultsJSON represents benchmark results in JSON format
type BenchmarkResultsJSON struct {
	TotalRequests      int64                  `json:"total_requests"`
	TotalErrors        int64                  `json:"total_errors"`
	Duration           string                 `json:"duration"`
	AverageLatency     string                 `json:"average_latency"`
	MinLatency         string                 `json:"min_latency"`
	MaxLatency         string                 `json:"max_latency"`
	LatencyStdDev      string                 `json:"latency_stddev"`
	RequestsPerSec     float64                `json:"requests_per_sec"`
	StatusCodes        map[int]int64          `json:"status_codes"`
	LatencyPercentiles map[string]string      `json:"latency_percentiles"` // Changed from map[float64]string to map[string]string
	TotalBytes         int64                  `json:"total_bytes"`
	EndpointStats      map[string]interface{} `json:"endpoint_stats,omitempty"`
}

// Server represents the API server
type Server struct {
	taskManager *TaskManager
}

// NewServer creates a new API server
func NewServer() *Server {
	return &Server{
		taskManager: NewTaskManager(),
	}
}

// generateTaskID generates a unique task ID
func generateTaskID() string {
	return fmt.Sprintf("task_%d", time.Now().UnixNano())
}

// handleBenchmark handles POST /api/v1/benchmark
func (s *Server) handleBenchmark(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req BenchmarkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// Validate request
	if req.Curl == "" && req.URL == "" {
		http.Error(w, "Either 'url' or 'curl' must be provided", http.StatusBadRequest)
		return
	}

	// Set defaults
	if req.Connections == 0 {
		req.Connections = 10
	}
	if req.Duration == "" {
		req.Duration = "10s"
	}
	if req.Threads == 0 {
		req.Threads = 2
	}
	if req.Timeout == "" {
		req.Timeout = "30s"
	}
	if req.Method == "" {
		req.Method = "GET"
	}

	// Parse duration and timeout
	duration, err := time.ParseDuration(req.Duration)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid duration format: %v", err), http.StatusBadRequest)
		return
	}

	timeout, err := time.ParseDuration(req.Timeout)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid timeout format: %v", err), http.StatusBadRequest)
		return
	}

	// Create config
	cfg := config.Config{
		Connections: req.Connections,
		Duration:    duration,
		Threads:     req.Threads,
		Rate:        req.Rate,
		Timeout:     timeout,
		CurlCommand: req.Curl,
		Method:      req.Method,
		Body:        req.Body,
		ContentType: req.ContentType,
		UseNetHTTP:  req.UseNetHTTP,
	}

	// Convert headers
	headers := make([]string, 0, len(req.Headers))
	for k, v := range req.Headers {
		headers = append(headers, fmt.Sprintf("%s: %s", k, v))
	}
	cfg.Headers = headers

	// Build HTTP request
	var httpReq *http.Request
	if req.Curl != "" {
		httpReq, err = parser.ParseCurl(req.Curl)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse curl command: %v", err), http.StatusBadRequest)
			return
		}
	} else {
		targetURL := req.URL
		if !strings.HasPrefix(targetURL, "http://") && !strings.HasPrefix(targetURL, "https://") {
			targetURL = "http://" + targetURL
		}

		parsedURL, err := url.Parse(targetURL)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid URL: %v", err), http.StatusBadRequest)
			return
		}

		httpReq, err = parser.BuildRequest(cfg, parsedURL)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to build request: %v", err), http.StatusBadRequest)
			return
		}
	}

	// Validate config
	if err := cfg.Validate(); err != nil {
		http.Error(w, fmt.Sprintf("Invalid configuration: %v", err), http.StatusBadRequest)
		return
	}

	// Create task
	taskID := generateTaskID()
	configMap := map[string]interface{}{
		"url":         req.URL,
		"curl":        req.Curl,
		"connections": req.Connections,
		"duration":    req.Duration,
		"threads":     req.Threads,
		"rate":        req.Rate,
		"timeout":     req.Timeout,
		"method":      req.Method,
		"headers":     req.Headers,
		"use_nethttp": req.UseNetHTTP,
	}
	if req.Extra != nil {
		for k, v := range req.Extra {
			configMap[k] = v
		}
	}

	s.taskManager.CreateTask(taskID, configMap)

	// Create benchmark runner
	bench := benchmark.New(cfg, httpReq)

	// Run task asynchronously with a new context (not tied to HTTP request)
	// This ensures the benchmark continues even after the HTTP request completes
	benchCtx := context.Background()
	s.taskManager.RunTask(benchCtx, taskID, func(ctx context.Context) (*stats.Results, error) {
		results, err := bench.Run(ctx)
		return results, err
	})

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(BenchmarkResponse{
		TaskID:  taskID,
		Status:  "accepted",
		Message: "Benchmark task created and started",
		Config:  configMap,
	})
}

// handleBatch handles POST /api/v1/batch
func (s *Server) handleBatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req BatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	if len(req.Tests) == 0 {
		http.Error(w, "No tests provided", http.StatusBadRequest)
		return
	}

	// Create batch task
	taskID := generateTaskID()
	configMap := map[string]interface{}{
		"type":  "batch",
		"tests": len(req.Tests),
	}

	s.taskManager.CreateTask(taskID, configMap)

	// Run batch tests asynchronously with a new context (not tied to HTTP request)
	benchCtx := context.Background()
	s.taskManager.RunTask(benchCtx, taskID, func(ctx context.Context) (*stats.Results, error) {
		// For batch tests, we'll run them sequentially and aggregate results
		// This is a simplified implementation
		// In a full implementation, you might want to use the batch executor
		return nil, fmt.Errorf("batch execution not yet implemented in API")
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(BenchmarkResponse{
		TaskID:  taskID,
		Status:  "accepted",
		Message: "Batch test task created",
		Config:  configMap,
	})
}

// handleStatus handles GET /api/v1/status/:id
func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract task ID from path
	path := r.URL.Path
	if !strings.HasPrefix(path, "/api/v1/status/") {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	taskID := strings.TrimPrefix(path, "/api/v1/status/")
	if taskID == "" {
		http.Error(w, "Task ID required", http.StatusBadRequest)
		return
	}

	task, exists := s.taskManager.GetTask(taskID)
	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(TaskStatusResponse{
		ID:          task.ID,
		Status:      string(task.Status),
		CreatedAt:   task.CreatedAt,
		StartedAt:   task.StartedAt,
		CompletedAt: task.CompletedAt,
		Error:       task.Error,
		Config:      task.Config,
	})
}

// handleResults handles GET /api/v1/results/:id
func (s *Server) handleResults(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract task ID from path
	path := r.URL.Path
	if !strings.HasPrefix(path, "/api/v1/results/") {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	taskID := strings.TrimPrefix(path, "/api/v1/results/")
	if taskID == "" {
		http.Error(w, "Task ID required", http.StatusBadRequest)
		return
	}

	task, exists := s.taskManager.GetTask(taskID)
	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	response := TaskResultsResponse{
		ID:          task.ID,
		Status:      string(task.Status),
		CreatedAt:   task.CreatedAt,
		StartedAt:   task.StartedAt,
		CompletedAt: task.CompletedAt,
		Error:       task.Error,
	}

	if task.Results != nil {
		response.Results = convertResultsToJSON(task.Results)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// convertResultsToJSON converts stats.Results to JSON format
func convertResultsToJSON(results *stats.Results) *BenchmarkResultsJSON {
	if results == nil {
		return nil
	}

	avgLatency := results.GetAverageLatency()
	minLatency := results.GetMinLatency()
	maxLatency := results.GetMaxLatency()
	stdDev := results.GetLatencyStdDev()
	percentiles := results.GetLatencyPercentiles()

	// Calculate requests per second
	var requestsPerSec float64
	if results.Duration > 0 {
		requestsPerSec = float64(results.TotalRequests) / results.Duration.Seconds()
	}

	// Convert percentiles to string format (JSON requires string keys)
	percentilesMap := make(map[string]string)
	for p, d := range percentiles {
		percentilesMap[fmt.Sprintf("p%.0f", p)] = formatDuration(d)
	}

	// Convert endpoint stats
	endpointStatsMap := make(map[string]interface{})
	for url, epStats := range results.GetEndpointStats() {
		epAvgLatency := epStats.GetAverageLatency()
		var epRequestsPerSec float64
		if results.Duration > 0 {
			epRequestsPerSec = float64(epStats.Requests) / results.Duration.Seconds()
		}

		endpointStatsMap[url] = map[string]interface{}{
			"requests":         epStats.Requests,
			"errors":           epStats.Errors,
			"requests_per_sec": epRequestsPerSec,
			"average_latency":  formatDuration(epAvgLatency),
			"min_latency":      formatDuration(epStats.MinLatency),
			"max_latency":      formatDuration(epStats.MaxLatency),
			"status_codes":     epStats.StatusCodes,
			"total_bytes":      epStats.TotalBytes,
		}
	}

	return &BenchmarkResultsJSON{
		TotalRequests:      results.TotalRequests,
		TotalErrors:        results.TotalErrors,
		Duration:           formatDuration(results.Duration),
		AverageLatency:     formatDuration(avgLatency),
		MinLatency:         formatDuration(minLatency),
		MaxLatency:         formatDuration(maxLatency),
		LatencyStdDev:      formatDuration(stdDev),
		RequestsPerSec:     requestsPerSec,
		StatusCodes:        results.GetStatusCodes(),
		LatencyPercentiles: percentilesMap,
		TotalBytes:         results.GetTotalBytes(),
		EndpointStats:      endpointStatsMap,
	}
}

// formatDuration formats a duration as a string
func formatDuration(d time.Duration) string {
	if d == 0 {
		return "0ms"
	}

	if d < time.Microsecond {
		return fmt.Sprintf("%dns", d.Nanoseconds())
	}
	if d < time.Millisecond {
		return fmt.Sprintf("%.2fμs", float64(d)/float64(time.Microsecond))
	}
	if d < time.Second {
		return fmt.Sprintf("%.2fms", float64(d)/float64(time.Millisecond))
	}
	return fmt.Sprintf("%.2fs", d.Seconds())
}

package mcp

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/antlabs/gurl/internal/benchmark"
	"github.com/antlabs/gurl/internal/config"
	"github.com/antlabs/gurl/internal/parser"
	"github.com/antlabs/gurl/internal/stats"
	"github.com/mark3labs/mcp-go/mcp"
)

// handleBenchmark handles the gurl.benchmark tool
func (s *Server) handleBenchmark(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	defer func() {
		if r := recover(); r != nil {
			Logger.Printf("Panic in handleBenchmark: %v", r)
		}
	}()

	// Parse arguments
	connections := mcp.ParseInt(req, "connections", 10)
	durationStr := mcp.ParseString(req, "duration", "10s")
	threads := mcp.ParseInt(req, "threads", 2)
	rate := mcp.ParseInt(req, "rate", 0)
	requests := mcp.ParseInt(req, "requests", 0)
	timeoutStr := mcp.ParseString(req, "timeout", "30s")
	curlCommand := mcp.ParseString(req, "curl", "")
	targetURL := mcp.ParseString(req, "url", "")
	method := mcp.ParseString(req, "method", "GET")
	headers := mcp.ParseStringMap(req, "headers", map[string]any{})
	body := mcp.ParseString(req, "body", "")
	contentType := mcp.ParseString(req, "content_type", "")
	verbose := mcp.ParseBoolean(req, "verbose", false)
	printLatency := mcp.ParseBoolean(req, "latency", false)
	useNetHTTP := mcp.ParseBoolean(req, "use_nethttp", false)

	// Parse duration and timeout
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		Logger.Printf("Invalid duration format: %v", err)
		return nil, fmt.Errorf("invalid duration format: %w", err)
	}

	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		Logger.Printf("Invalid timeout format: %v", err)
		return nil, fmt.Errorf("invalid timeout format: %w", err)
	}

	// Create config
	cfg := config.Config{
		Connections:  connections,
		Duration:     duration,
		Threads:      threads,
		Rate:         rate,
		Requests:     int64(requests),
		Timeout:      timeout,
		CurlCommand:  curlCommand,
		Method:       method,
		Headers:      []string{}, // We'll add headers individually
		Body:         body,
		ContentType:  contentType,
		Verbose:      verbose,
		PrintLatency: printLatency,
		UseNetHTTP:   useNetHTTP,
	}

	var httpReq *http.Request
	// Handle curl command or build request from parameters
	if curlCommand != "" {
		Logger.Printf("Parsing curl command: %s", curlCommand)
		httpReq, err = parser.ParseCurl(curlCommand)
		if err != nil {
			Logger.Printf("Failed to parse curl command: %v", err)
			return nil, fmt.Errorf("failed to parse curl command: %w", err)
		}
	} else {
		if targetURL == "" {
			Logger.Printf("URL is required when not using curl command")
			return nil, fmt.Errorf("url is required when not using curl command")
		}

		// Validate and format URL
		if !strings.HasPrefix(targetURL, "http://") && !strings.HasPrefix(targetURL, "https://") {
			targetURL = "http://" + targetURL
		}

		parsedURL, err := url.Parse(targetURL)
		if err != nil {
			Logger.Printf("Invalid URL: %v", err)
			return nil, fmt.Errorf("invalid URL: %w", err)
		}

		httpReq, err = parser.BuildRequest(cfg, parsedURL)
		if err != nil {
			Logger.Printf("Failed to build request: %v", err)
			return nil, fmt.Errorf("failed to build request: %w", err)
		}
	}

	// Add headers
	for k, v := range headers {
		if vs, ok := v.(string); ok {
			httpReq.Header.Set(k, vs)
		}
	}

	// Validate config
	if err := cfg.Validate(); err != nil {
		Logger.Printf("Invalid configuration: %v", err)
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	// Create and run benchmark
	Logger.Printf("Starting benchmark with %d connections, %d threads, duration %s", connections, threads, durationStr)
	bench := benchmark.New(cfg, httpReq)

	// Create a new context for the benchmark
	benchCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	results, err := bench.Run(benchCtx)
	if err != nil {
		Logger.Printf("Benchmark failed: %v", err)
		return nil, fmt.Errorf("benchmark failed: %w", err)
	}

	Logger.Printf("Benchmark completed successfully with %d requests", results.TotalRequests)

	// Format results
	result := formatBenchmarkResults(results, cfg)

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Type: "text",
				Text: result,
			},
		},
	}, nil
}

// formatBenchmarkResults formats the benchmark results similar to wrk output
func formatBenchmarkResults(results *stats.Results, cfg config.Config) string {
	var result strings.Builder

	// Header
	result.WriteString(fmt.Sprintf("Running %s test\n", cfg.Duration))
	result.WriteString(fmt.Sprintf("  %d threads and %d connections\n", cfg.Threads, cfg.Connections))

	// Thread Stats
	result.WriteString("  Thread Stats   Avg      Stdev     Max   +/- Stdev\n")

	// Calculate latency stats
	latencies := results.GetLatencies()
	if len(latencies) > 0 {
		var total time.Duration
		for _, lat := range latencies {
			total += lat
		}
		avg := total / time.Duration(len(latencies))

		// Calculate standard deviation
		var sumSquares float64
		for _, lat := range latencies {
			diff := float64(lat - avg)
			sumSquares += diff * diff
		}
		variance := sumSquares / float64(len(latencies))
		stdev := time.Duration(variance)

		// Find max
		max := time.Duration(0)
		for _, lat := range latencies {
			if lat > max {
				max = lat
			}
		}

		result.WriteString(fmt.Sprintf("    Latency   %8s %8s %8s %8s\n",
			formatDuration(avg),
			formatDuration(stdev),
			formatDuration(max),
			"N/A"))
	}

	// Calculate RPS
	rps := float64(results.TotalRequests) / cfg.Duration.Seconds()
	result.WriteString(fmt.Sprintf("    Req/Sec   %8.2f %8s %8s %8s\n", rps, "N/A", "N/A", "N/A"))

	// Latency Distribution
	if cfg.PrintLatency && len(latencies) > 0 {
		result.WriteString("  Latency Distribution\n")
		percentiles := []float64{50, 75, 90, 99}
		for _, p := range percentiles {
			idx := int(float64(len(latencies)-1) * p / 100.0)
			if idx >= 0 && idx < len(latencies) {
				result.WriteString(fmt.Sprintf("     %2.0f%%   %s\n", p, formatDuration(latencies[idx])))
			}
		}
	}

	// Summary
	result.WriteString(fmt.Sprintf("  %d requests in %s\n", results.TotalRequests, cfg.Duration))

	// Status code distribution
	statusCodes := results.GetStatusCodes()
	if len(statusCodes) > 0 {
		result.WriteString("  Status code distribution:\n")
		for code, count := range statusCodes {
			percentage := float64(count) / float64(results.TotalRequests) * 100
			result.WriteString(fmt.Sprintf("    [%d] %d responses (%.1f%%)\n", code, count, percentage))
		}
	}

	result.WriteString(fmt.Sprintf("Requests/sec: %8.2f\n", rps))

	return result.String()
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

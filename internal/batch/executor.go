package batch

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/antlabs/gurl/internal/benchmark"
	"github.com/antlabs/gurl/internal/config"
	"github.com/antlabs/gurl/internal/parser"
	"github.com/antlabs/gurl/internal/stats"
)

// TestResult represents the result of a single batch test
type TestResult struct {
	Name      string
	Config    *config.Config
	Stats     *stats.Results
	Error     error
	StartTime time.Time
	EndTime   time.Time
	Duration  time.Duration
}

// BatchResult represents the result of a batch test run
type BatchResult struct {
	Tests       []TestResult
	TotalTime   time.Duration
	SuccessRate float64
	StartTime   time.Time
	EndTime     time.Time
}

// Executor handles batch test execution
type Executor struct {
	maxConcurrency int
	verbose        bool
}

// NewExecutor creates a new batch executor
func NewExecutor(maxConcurrency int, verbose bool) *Executor {
	if maxConcurrency <= 0 {
		maxConcurrency = 1
	}
	return &Executor{
		maxConcurrency: maxConcurrency,
		verbose:        verbose,
	}
}

// Execute runs all tests in the batch configuration
func (e *Executor) Execute(ctx context.Context, batchConfig *config.BatchConfig, defaults *config.Config) (*BatchResult, error) {
	if err := batchConfig.Validate(); err != nil {
		return nil, fmt.Errorf("invalid batch configuration: %v", err)
	}

	startTime := time.Now()

	// Create semaphore to limit concurrency
	sem := make(chan struct{}, e.maxConcurrency)
	var wg sync.WaitGroup
	results := make([]TestResult, len(batchConfig.Tests))

	if e.verbose {
		fmt.Printf("Starting batch test with %d tests (max concurrency: %d)\n", len(batchConfig.Tests), e.maxConcurrency)
	}

	// Execute tests concurrently
	for i, test := range batchConfig.Tests {
		wg.Add(1)
		go func(index int, batchTest config.BatchTest) {
			defer wg.Done()

			// Acquire semaphore
			sem <- struct{}{}
			defer func() { <-sem }()

			result := e.executeTest(ctx, &batchTest, defaults)
			results[index] = result

			if e.verbose {
				if result.Error != nil {
					fmt.Printf("Test '%s' failed: %v\n", result.Name, result.Error)
				} else {
					fmt.Printf("Test '%s' completed in %v\n", result.Name, result.Duration)
				}
			}
		}(i, test)
	}

	// Wait for all tests to complete
	wg.Wait()
	endTime := time.Now()

	// Calculate success rate
	successCount := 0
	for _, result := range results {
		// A test is considered successful only if there is no top-level error
		// and no per-request errors recorded in stats (e.g. assertion failures).
		if result.Error == nil {
			if result.Stats == nil || len(result.Stats.GetErrors()) == 0 {
				successCount++
			}
		}
	}
	successRate := float64(successCount) / float64(len(results)) * 100

	return &BatchResult{
		Tests:       results,
		TotalTime:   endTime.Sub(startTime),
		SuccessRate: successRate,
		StartTime:   startTime,
		EndTime:     endTime,
	}, nil
}

// executeTest runs a single test
func (e *Executor) executeTest(ctx context.Context, batchTest *config.BatchTest, defaults *config.Config) TestResult {
	startTime := time.Now()

	result := TestResult{
		Name:      batchTest.Name,
		StartTime: startTime,
	}

	// Convert batch test to config
	cfg, err := batchTest.ToConfig(defaults)
	if err != nil {
		result.Error = fmt.Errorf("failed to create config: %v", err)
		result.EndTime = time.Now()
		result.Duration = result.EndTime.Sub(result.StartTime)
		return result
	}
	result.Config = cfg

	// Parse curl command if provided and create http.Request
	var req *http.Request
	if cfg.CurlCommand != "" {
		parsedReq, err := parser.ParseCurl(cfg.CurlCommand)
		if err != nil {
			result.Error = fmt.Errorf("failed to parse curl command: %v", err)
			result.EndTime = time.Now()
			result.Duration = result.EndTime.Sub(result.StartTime)
			return result
		}
		req = parsedReq

		// Update config with parsed request
		cfg.Method = req.Method
		cfg.Headers = make([]string, 0, len(req.Header))
		for key, values := range req.Header {
			for _, value := range values {
				cfg.Headers = append(cfg.Headers, fmt.Sprintf("%s: %s", key, value))
			}
		}
		if req.Body != nil {
			body, err := io.ReadAll(req.Body)
			if err != nil {
				result.Error = fmt.Errorf("failed to read request body: %v", err)
				result.EndTime = time.Now()
				result.Duration = result.EndTime.Sub(result.StartTime)
				return result
			}
			cfg.Body = string(body)
		}
	} else {
		// TODO raw http request
		// TODO:
		// Create a basic request if no curl command provided
		var err error
		req, err = http.NewRequest(cfg.Method, "http://example.com", nil)
		if err != nil {
			result.Error = fmt.Errorf("failed to create request: %v", err)
			result.EndTime = time.Now()
			result.Duration = result.EndTime.Sub(result.StartTime)
			return result
		}
	}

	// Validate final config
	if err := cfg.Validate(); err != nil {
		result.Error = fmt.Errorf("invalid configuration: %v", err)
		result.EndTime = time.Now()
		result.Duration = result.EndTime.Sub(result.StartTime)
		return result
	}

	// Create and run benchmark
	bench := benchmark.New(*cfg, req)

	// Run the benchmark
	benchStats, err := bench.Run(ctx)
	if err != nil {
		result.Error = fmt.Errorf("benchmark failed: %v", err)
		result.EndTime = time.Now()
		result.Duration = result.EndTime.Sub(result.StartTime)
		return result
	}

	result.Stats = benchStats
	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)

	return result
}

// ExecuteSequential runs tests sequentially (for debugging or when concurrency is not desired)
func (e *Executor) ExecuteSequential(ctx context.Context, batchConfig *config.BatchConfig, defaults *config.Config) (*BatchResult, error) {
	if err := batchConfig.Validate(); err != nil {
		return nil, fmt.Errorf("invalid batch configuration: %v", err)
	}

	startTime := time.Now()
	results := make([]TestResult, 0, len(batchConfig.Tests))

	if e.verbose {
		fmt.Printf("Starting sequential batch test with %d tests\n", len(batchConfig.Tests))
	}

	for i, test := range batchConfig.Tests {
		if e.verbose {
			fmt.Printf("Running test %d/%d: %s\n", i+1, len(batchConfig.Tests), test.Name)
		}

		result := e.executeTest(ctx, &test, defaults)
		results = append(results, result)

		if result.Error != nil && e.verbose {
			fmt.Printf("Test '%s' failed: %v\n", result.Name, result.Error)
		}
	}

	endTime := time.Now()

	// Calculate success rate
	successCount := 0
	for _, result := range results {
		// A test is considered successful only if there is no top-level error
		// and no per-request errors recorded in stats (e.g. assertion failures).
		if result.Error == nil {
			if result.Stats == nil || len(result.Stats.GetErrors()) == 0 {
				successCount++
			}
		}
	}
	successRate := float64(successCount) / float64(len(results)) * 100

	return &BatchResult{
		Tests:       results,
		TotalTime:   endTime.Sub(startTime),
		SuccessRate: successRate,
		StartTime:   startTime,
		EndTime:     endTime,
	}, nil
}

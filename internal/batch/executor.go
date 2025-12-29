package batch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/antlabs/gurl/internal/benchmark"
	"github.com/antlabs/gurl/internal/config"
	"github.com/antlabs/gurl/internal/parser"
	"github.com/antlabs/gurl/internal/stats"
)

// TestResult represents the result of a single batch test
type TestResult struct {
	Name           string
	Config         *config.Config
	Stats          *stats.Results
	Error          error
	StartTime      time.Time
	EndTime        time.Time
	Duration       time.Duration
	SampleResponse []byte // Store a sample response for jsondiff comparison
}

// BatchResult represents the result of a batch test run
type BatchResult struct {
	Tests           []TestResult
	TotalTime       time.Duration
	SuccessRate     float64
	StartTime       time.Time
	EndTime         time.Time
	JSONDiffResults []JSONDiffResult // Results of JSON comparisons
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

	batchResult := &BatchResult{
		Tests:       results,
		TotalTime:   endTime.Sub(startTime),
		SuccessRate: successRate,
		StartTime:   startTime,
		EndTime:     endTime,
	}

	// Perform jsondiff comparisons if configured
	if batchConfig.JSONDiff != nil && len(batchConfig.JSONDiff.Pairs) > 0 {
		batchResult.JSONDiffResults = e.performJSONDiffComparisons(batchConfig.JSONDiff, results)
	}

	return batchResult, nil
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

	// Capture a sample response before running the benchmark (for jsondiff)
	sampleResp, err := captureSampleResponse(req)
	if err == nil {
		result.SampleResponse = sampleResp
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

	// Save result to file if output_file is specified
	if batchTest.OutputFile != "" {
		if err := e.saveTestResultToFile(&result, batchTest.OutputFile, batchTest.OutputResponseOnly); err != nil {
			if e.verbose {
				fmt.Printf("Warning: failed to save result to file '%s': %v\n", batchTest.OutputFile, err)
			}
		} else if e.verbose {
			fmt.Printf("Test result saved to: %s\n", batchTest.OutputFile)
		}
	}

	return result
}

// saveTestResultToFile saves a test result to a JSON file
func (e *Executor) saveTestResultToFile(result *TestResult, filename string, responseOnly bool) error {
	// If response-only mode, just save the sample response
	if responseOnly {
		if len(result.SampleResponse) == 0 {
			return fmt.Errorf("no sample response available")
		}

		// Try to format as JSON if it's valid JSON
		var jsonData interface{}
		if err := json.Unmarshal(result.SampleResponse, &jsonData); err == nil {
			// It's valid JSON, format it with indentation
			formattedJSON, err := json.MarshalIndent(jsonData, "", "  ")
			if err == nil {
				if err := os.WriteFile(filename, formattedJSON, 0644); err != nil {
					return fmt.Errorf("failed to write file: %v", err)
				}
				return nil
			}
		}

		// Not valid JSON or formatting failed, write raw response
		if err := os.WriteFile(filename, result.SampleResponse, 0644); err != nil {
			return fmt.Errorf("failed to write file: %v", err)
		}
		return nil
	}

	// Create output structure (full statistics mode)
	output := map[string]interface{}{
		"name":       result.Name,
		"start_time": result.StartTime.Format(time.RFC3339),
		"end_time":   result.EndTime.Format(time.RFC3339),
		"duration":   result.Duration.String(),
	}

	if result.Error != nil {
		output["error"] = result.Error.Error()
		output["status"] = "failed"
	} else {
		output["status"] = "success"
	}

	// Add statistics if available
	if result.Stats != nil {
		stats := map[string]interface{}{
			"total_requests": result.Stats.TotalRequests,
			"total_errors":   result.Stats.TotalErrors,
			"duration":       result.Stats.Duration.String(),
		}

		// Calculate RPS
		if result.Stats.Duration > 0 {
			rps := float64(result.Stats.TotalRequests) / result.Stats.Duration.Seconds()
			stats["rps"] = fmt.Sprintf("%.2f", rps)
		}

		// Add latency statistics
		avgLatency := result.Stats.GetAverageLatency()
		if avgLatency > 0 {
			stats["avg_latency_ms"] = avgLatency.Milliseconds()
		}

		minLatency := result.Stats.GetMinLatency()
		if minLatency > 0 {
			stats["min_latency_ms"] = minLatency.Milliseconds()
		}

		maxLatency := result.Stats.GetMaxLatency()
		if maxLatency > 0 {
			stats["max_latency_ms"] = maxLatency.Milliseconds()
		}

		// Add percentiles
		percentiles := result.Stats.GetLatencyPercentiles()
		if len(percentiles) > 0 {
			pctMap := make(map[string]int64)
			for p, latency := range percentiles {
				pctMap[fmt.Sprintf("p%d", int(p*100))] = latency.Milliseconds()
			}
			stats["percentiles"] = pctMap
		}

		// Add status codes
		statusCodes := result.Stats.GetStatusCodes()
		if len(statusCodes) > 0 {
			stats["status_codes"] = statusCodes
		}

		output["statistics"] = stats

		// Add errors if any
		errs := result.Stats.GetErrors()
		if len(errs) > 0 {
			errorMsgs := make([]string, 0, len(errs))
			for _, err := range errs {
				if err != nil {
					errorMsgs = append(errorMsgs, err.Error())
				}
			}
			if len(errorMsgs) > 0 {
				output["errors"] = errorMsgs
			}
		}
	}

	// Add sample response if available
	if len(result.SampleResponse) > 0 {
		output["sample_response"] = string(result.SampleResponse)
	}

	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal result: %v", err)
	}

	// Write to file
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}

// captureSampleResponse captures a single response from the endpoint for jsondiff comparison
func captureSampleResponse(req *http.Request) ([]byte, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Clone the request to avoid modifying the original
	clonedReq := req.Clone(context.Background())

	resp, err := client.Do(clonedReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// performJSONDiffComparisons performs JSON diff comparisons based on configuration
func (e *Executor) performJSONDiffComparisons(jsonDiffConfig *config.JSONDiffConfig, results []TestResult) []JSONDiffResult {
	diffResults := make([]JSONDiffResult, 0, len(jsonDiffConfig.Pairs))

	for _, pair := range jsonDiffConfig.Pairs {
		// Find base and target test results
		var baseResult, targetResult *TestResult
		for i := range results {
			if results[i].Name == pair.Base {
				baseResult = &results[i]
			}
			if results[i].Name == pair.Target {
				targetResult = &results[i]
			}
		}

		// Check if both tests were found
		if baseResult == nil {
			diffResults = append(diffResults, JSONDiffResult{
				PairName:   pair.Name,
				BaseName:   pair.Base,
				TargetName: pair.Target,
				Passed:     false,
				Message:    fmt.Sprintf("base test '%s' not found", pair.Base),
			})
			continue
		}

		if targetResult == nil {
			diffResults = append(diffResults, JSONDiffResult{
				PairName:   pair.Name,
				BaseName:   pair.Base,
				TargetName: pair.Target,
				Passed:     false,
				Message:    fmt.Sprintf("target test '%s' not found", pair.Target),
			})
			continue
		}

		// Check if both tests have sample responses
		if len(baseResult.SampleResponse) == 0 {
			diffResults = append(diffResults, JSONDiffResult{
				PairName:   pair.Name,
				BaseName:   pair.Base,
				TargetName: pair.Target,
				Passed:     false,
				Message:    fmt.Sprintf("base test '%s' has no sample response", pair.Base),
			})
			continue
		}

		if len(targetResult.SampleResponse) == 0 {
			diffResults = append(diffResults, JSONDiffResult{
				PairName:   pair.Name,
				BaseName:   pair.Base,
				TargetName: pair.Target,
				Passed:     false,
				Message:    fmt.Sprintf("target test '%s' has no sample response", pair.Target),
			})
			continue
		}

		// Perform the comparison
		pairName := pair.Name
		if pairName == "" {
			pairName = fmt.Sprintf("%s vs %s", pair.Base, pair.Target)
		}

		// Determine which fields to compare
		baseField := pair.BaseField
		targetField := pair.TargetField

		// If base_field and target_field are not specified, fall back to compare_field
		if baseField == "" && targetField == "" && pair.CompareField != "" {
			baseField = pair.CompareField
			targetField = pair.CompareField
		}

		result := CompareJSONResponsesWithFields(
			pair.Base,
			pair.Target,
			baseResult.SampleResponse,
			targetResult.SampleResponse,
			baseField,
			targetField,
			pair.IgnoreFields,
		)
		result.PairName = pairName

		diffResults = append(diffResults, result)

		if e.verbose {
			if result.Passed {
				fmt.Printf("JSONDiff '%s': PASSED - %s\n", pairName, result.Message)
			} else {
				fmt.Printf("JSONDiff '%s': FAILED - %s\n", pairName, result.Message)
			}
		}
	}

	return diffResults
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

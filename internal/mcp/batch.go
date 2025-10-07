package mcp

import (
	"context"
	"fmt"
	"time"

	"github.com/antlabs/gurl/internal/batch"
	"github.com/antlabs/gurl/internal/config"
	"github.com/mark3labs/mcp-go/mcp"
)

// handleBatchTest handles the gurl.batch_test tool
func (s *Server) handleBatchTest(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	defer func() {
		if r := recover(); r != nil {
			Logger.Printf("Panic in handleBatchTest: %v", r)
		}
	}()
	
	// Parse arguments
	batchConfigPath := mcp.ParseString(req, "config", "")
	inlineTests := mcp.ParseArgument(req, "tests", []any{})
	verbose := mcp.ParseBoolean(req, "verbose", false)
	maxConcurrency := mcp.ParseInt(req, "concurrency", 3)
	
	// Check if either config file or inline tests are provided
	if batchConfigPath == "" && len(inlineTests.([]any)) == 0 {
		Logger.Printf("Either config file path or inline tests are required")
		return nil, fmt.Errorf("either config file path or inline tests are required")
	}
	
	// Create default config
	defaults := &config.Config{
		Connections: 10,
		Duration:    10 * time.Second,
		Threads:     2,
		Rate:        0,
		Timeout:     30 * time.Second,
	}

	// Create batch executor
	executor := batch.NewExecutor(maxConcurrency, verbose)

	// Create context
	batchCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var batchConfig *config.BatchConfig
	var err error

	// Load batch configuration from file or parse inline tests
	if batchConfigPath != "" {
		Logger.Printf("Loading batch configuration from: %s", batchConfigPath)
		batchConfig, err = config.LoadBatchConfig(batchConfigPath)
		if err != nil {
			Logger.Printf("Failed to load batch config: %v", err)
			return nil, fmt.Errorf("failed to load batch config: %w", err)
		}
	} else {
		// Parse inline tests
		Logger.Printf("Parsing inline batch tests")
		batchConfig, err = parseInlineTests(inlineTests.([]any))
		if err != nil {
			Logger.Printf("Failed to parse inline tests: %v", err)
			return nil, fmt.Errorf("failed to parse inline tests: %w", err)
		}
	}

	// Execute batch tests
	Logger.Printf("Executing batch tests with max concurrency: %d", maxConcurrency)
	result, err := executor.Execute(batchCtx, batchConfig, defaults)
	if err != nil {
		Logger.Printf("Batch test failed: %v", err)
		return nil, fmt.Errorf("batch test failed: %w", err)
	}
	
	Logger.Printf("Batch tests completed successfully")

	// Generate report
	reporter := batch.NewReporter(verbose)
	report := reporter.GenerateReport(result)

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: report,
			},
		},
	}, nil
}

// parseInlineTests parses inline test definitions
func parseInlineTests(inlineTests []any) (*config.BatchConfig, error) {
	batchConfig := &config.BatchConfig{
		Version: "1.0",
		Tests:   make([]config.BatchTest, 0, len(inlineTests)),
	}

	for i, test := range inlineTests {
		testMap, ok := test.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("invalid test format at index %d", i)
		}

		batchTest := config.BatchTest{}

		// Parse test fields
		if name, ok := testMap["name"].(string); ok {
			batchTest.Name = name
		} else {
			batchTest.Name = fmt.Sprintf("Test %d", i+1)
		}

		if curl, ok := testMap["curl"].(string); ok {
			batchTest.Curl = curl
		}

		if connections, ok := testMap["connections"].(float64); ok {
			batchTest.Connections = int(connections)
		}

		if duration, ok := testMap["duration"].(string); ok {
			batchTest.Duration = duration
		}

		if threads, ok := testMap["threads"].(float64); ok {
			batchTest.Threads = int(threads)
		}

		if rate, ok := testMap["rate"].(float64); ok {
			batchTest.Rate = int(rate)
		}

		if timeout, ok := testMap["timeout"].(string); ok {
			batchTest.Timeout = timeout
		}

		if verbose, ok := testMap["verbose"].(bool); ok {
			batchTest.Verbose = verbose
		}

		if useNetHTTP, ok := testMap["use_nethttp"].(bool); ok {
			batchTest.UseNetHTTP = useNetHTTP
		}

		batchConfig.Tests = append(batchConfig.Tests, batchTest)
	}

	return batchConfig, nil
}
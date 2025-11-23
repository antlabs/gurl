package batch

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// Reporter handles batch test result reporting
type Reporter struct {
	verbose bool
}

// NewReporter creates a new batch reporter
func NewReporter(verbose bool) *Reporter {
	return &Reporter{
		verbose: verbose,
	}
}

// GenerateReport generates a detailed batch test report
func (r *Reporter) GenerateReport(result *BatchResult) string {
	var report strings.Builder

	// Header
	report.WriteString("=== Batch Test Report ===\n\n")

	// Summary
	report.WriteString(fmt.Sprintf("Total Tests: %d\n", len(result.Tests)))
	report.WriteString(fmt.Sprintf("Success Rate: %.2f%%\n", result.SuccessRate))
	report.WriteString(fmt.Sprintf("Total Time: %v\n", result.TotalTime))
	report.WriteString(fmt.Sprintf("Start Time: %s\n", result.StartTime.Format("2006-01-02 15:04:05")))
	report.WriteString(fmt.Sprintf("End Time: %s\n\n", result.EndTime.Format("2006-01-02 15:04:05")))

	// Test Results
	report.WriteString("=== Test Results ===\n\n")

	successCount := 0
	failedCount := 0

	for i, test := range result.Tests {
		report.WriteString(fmt.Sprintf("%d. %s\n", i+1, test.Name))
		report.WriteString(fmt.Sprintf("   Duration: %v\n", test.Duration))

		var statsErrorCount int
		if test.Stats != nil {
			statsErrorCount = len(test.Stats.GetErrors())
		}

		// A test is considered FAILED if there is a top-level error
		// or any per-request errors (e.g. assertion failures).
		if test.Error != nil || statsErrorCount > 0 {
			report.WriteString("   Status: FAILED\n")
			if test.Error != nil {
				report.WriteString(fmt.Sprintf("   Error: %v\n", test.Error))
			}
			if test.Stats != nil {
				report.WriteString(fmt.Sprintf("   Requests: %d\n", test.Stats.TotalRequests))
				rps := float64(test.Stats.TotalRequests) / test.Stats.Duration.Seconds()
				report.WriteString(fmt.Sprintf("   RPS: %.2f\n", rps))
				report.WriteString(fmt.Sprintf("   Avg Latency: %v\n", test.Stats.GetAverageLatency()))
				if statsErrorCount > 0 {
					report.WriteString(fmt.Sprintf("   Errors: %d\n", statsErrorCount))
					// Print a few sample errors (typically assertion failures) in a
					// compact, readable block.
					report.WriteString("   Assertion Errors:\n")
					const maxSampleErrors = 5
					for idx, err := range test.Stats.GetErrors() {
						if idx >= maxSampleErrors {
							break
						}
						report.WriteString(fmt.Sprintf("     - #%d: %v\n", idx+1, err))
					}
				}
			}
			failedCount++
		} else {
			report.WriteString("   Status: SUCCESS\n")
			if test.Stats != nil {
				report.WriteString(fmt.Sprintf("   Requests: %d\n", test.Stats.TotalRequests))
				rps := float64(test.Stats.TotalRequests) / test.Stats.Duration.Seconds()
				report.WriteString(fmt.Sprintf("   RPS: %.2f\n", rps))
				report.WriteString(fmt.Sprintf("   Avg Latency: %v\n", test.Stats.GetAverageLatency()))
			}
			successCount++
		}

		if r.verbose && test.Config != nil {
			report.WriteString(fmt.Sprintf("   Config: c=%d, t=%d, d=%v\n",
				test.Config.Connections, test.Config.Threads, test.Config.Duration))
		}

		report.WriteString("\n")
	}

	// Performance Summary
	if successCount > 0 {
		report.WriteString("=== Performance Summary ===\n\n")

		var totalRequests int64
		var totalRPS float64
		var latencies []time.Duration

		for _, test := range result.Tests {
			if test.Error == nil && test.Stats != nil {
				totalRequests += test.Stats.TotalRequests
				rps := float64(test.Stats.TotalRequests) / test.Stats.Duration.Seconds()
				totalRPS += rps
				latencies = append(latencies, test.Stats.GetAverageLatency())
			}
		}

		report.WriteString(fmt.Sprintf("Total Requests: %d\n", totalRequests))
		report.WriteString(fmt.Sprintf("Combined RPS: %.2f\n", totalRPS))

		if len(latencies) > 0 {
			sort.Slice(latencies, func(i, j int) bool {
				return latencies[i] < latencies[j]
			})

			avgLatency := r.calculateAverageLatency(latencies)
			medianLatency := latencies[len(latencies)/2]
			minLatency := latencies[0]
			maxLatency := latencies[len(latencies)-1]

			report.WriteString("Latency Stats:\n")
			report.WriteString(fmt.Sprintf("  Average: %v\n", avgLatency))
			report.WriteString(fmt.Sprintf("  Median:  %v\n", medianLatency))
			report.WriteString(fmt.Sprintf("  Min:     %v\n", minLatency))
			report.WriteString(fmt.Sprintf("  Max:     %v\n", maxLatency))
		}
		report.WriteString("\n")
	}

	// Failed Tests Summary
	if failedCount > 0 {
		report.WriteString("=== Failed Tests ===\n\n")
		for _, test := range result.Tests {
			var statsErrorCount int
			if test.Stats != nil {
				statsErrorCount = len(test.Stats.GetErrors())
			}

			if test.Error != nil || statsErrorCount > 0 {
				if test.Error != nil {
					report.WriteString(fmt.Sprintf("- %s\n", test.Name))
					report.WriteString(fmt.Sprintf("    Top-level error: %v\n", test.Error))
				} else if statsErrorCount > 0 {
					// Show first assertion/response error as a summary
					errs := test.Stats.GetErrors()
					if len(errs) > 0 {
						report.WriteString(fmt.Sprintf("- %s\n", test.Name))
						report.WriteString(fmt.Sprintf("    First error: %v\n", errs[0]))
						report.WriteString(fmt.Sprintf("    Total errors: %d\n", statsErrorCount))
					}
				}
			}
		}
		report.WriteString("\n")
	}

	return report.String()
}

// GenerateCSVReport generates a CSV format report
func (r *Reporter) GenerateCSVReport(result *BatchResult) string {
	var csv strings.Builder

	// CSV Header
	csv.WriteString("Name,Status,Duration,Requests,RPS,AvgLatency,Errors,Error\n")

	// CSV Data
	for _, test := range result.Tests {
		status := "SUCCESS"
		requests := "0"
		rps := "0"
		avgLatency := "0"
		errors := "0"
		errorMsg := ""

		if test.Error != nil {
			status = "FAILED"
			errorMsg = strings.ReplaceAll(test.Error.Error(), ",", ";")
		} else if test.Stats != nil {
			requests = fmt.Sprintf("%d", test.Stats.TotalRequests)
			rpsVal := float64(test.Stats.TotalRequests) / test.Stats.Duration.Seconds()
			rps = fmt.Sprintf("%.2f", rpsVal)
			avgLatency = test.Stats.GetAverageLatency().String()
			errorCount := len(test.Stats.GetErrors())
			errors = fmt.Sprintf("%d", errorCount)
		}

		csv.WriteString(fmt.Sprintf("%s,%s,%v,%s,%s,%s,%s,%s\n",
			test.Name, status, test.Duration, requests, rps, avgLatency, errors, errorMsg))
	}

	return csv.String()
}

// GenerateJSONReport generates a JSON format report
func (r *Reporter) GenerateJSONReport(result *BatchResult) string {
	// Simple JSON generation without external dependencies
	var json strings.Builder

	json.WriteString("{\n")
	json.WriteString(fmt.Sprintf("  \"total_tests\": %d,\n", len(result.Tests)))
	json.WriteString(fmt.Sprintf("  \"success_rate\": %.2f,\n", result.SuccessRate))
	json.WriteString(fmt.Sprintf("  \"total_time\": \"%v\",\n", result.TotalTime))
	json.WriteString(fmt.Sprintf("  \"start_time\": \"%s\",\n", result.StartTime.Format(time.RFC3339)))
	json.WriteString(fmt.Sprintf("  \"end_time\": \"%s\",\n", result.EndTime.Format(time.RFC3339)))
	json.WriteString("  \"tests\": [\n")

	for i, test := range result.Tests {
		json.WriteString("    {\n")
		json.WriteString(fmt.Sprintf("      \"name\": \"%s\",\n", test.Name))
		json.WriteString(fmt.Sprintf("      \"duration\": \"%v\",\n", test.Duration))

		if test.Error != nil {
			json.WriteString("      \"status\": \"FAILED\",\n")
			errorMsg := strings.ReplaceAll(test.Error.Error(), "\"", "\\\"")
			json.WriteString(fmt.Sprintf("      \"error\": \"%s\"\n", errorMsg))
		} else {
			json.WriteString("      \"status\": \"SUCCESS\"")
			if test.Stats != nil {
				json.WriteString(",\n")
				json.WriteString(fmt.Sprintf("      \"requests\": %d,\n", test.Stats.TotalRequests))
				rpsVal := float64(test.Stats.TotalRequests) / test.Stats.Duration.Seconds()
				json.WriteString(fmt.Sprintf("      \"rps\": %.2f,\n", rpsVal))
				json.WriteString(fmt.Sprintf("      \"avg_latency\": \"%v\",\n", test.Stats.GetAverageLatency()))
				errorCount := len(test.Stats.GetErrors())
				json.WriteString(fmt.Sprintf("      \"errors\": %d\n", errorCount))
			} else {
				json.WriteString("\n")
			}
		}

		if i < len(result.Tests)-1 {
			json.WriteString("    },\n")
		} else {
			json.WriteString("    }\n")
		}
	}

	json.WriteString("  ]\n")
	json.WriteString("}\n")

	return json.String()
}

// calculateAverageLatency calculates the average of a slice of durations
func (r *Reporter) calculateAverageLatency(latencies []time.Duration) time.Duration {
	if len(latencies) == 0 {
		return 0
	}

	var total time.Duration
	for _, latency := range latencies {
		total += latency
	}

	return total / time.Duration(len(latencies))
}

// PrintSummary prints a quick summary to stdout
func (r *Reporter) PrintSummary(result *BatchResult) {
	successCount := 0
	for _, test := range result.Tests {
		if test.Error == nil {
			successCount++
		}
	}

	fmt.Printf("\n=== Batch Test Summary ===\n")
	fmt.Printf("Tests: %d/%d passed (%.1f%%)\n", successCount, len(result.Tests), result.SuccessRate)
	fmt.Printf("Total time: %v\n", result.TotalTime)

	if successCount < len(result.Tests) {
		fmt.Printf("\nFailed tests:\n")
		for _, test := range result.Tests {
			if test.Error != nil {
				fmt.Printf("  - %s: %v\n", test.Name, test.Error)
			}
		}
	}
	fmt.Println()
}

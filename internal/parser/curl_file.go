package parser

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// ParseCurlFile reads and parses multiple curl commands from a file
// Each line should contain one curl command
// Empty lines and lines starting with # are ignored
func ParseCurlFile(filePath string) ([]*http.Request, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open curl file: %w", err)
	}
	defer func() { _ = file.Close() }()

	var requests []*http.Request
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse the curl command
		req, err := ParseCurl(line)
		if err != nil {
			return nil, fmt.Errorf("failed to parse curl command at line %d: %w", lineNum, err)
		}

		requests = append(requests, req)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading curl file: %w", err)
	}

	if len(requests) == 0 {
		return nil, fmt.Errorf("no valid curl commands found in file")
	}

	return requests, nil
}

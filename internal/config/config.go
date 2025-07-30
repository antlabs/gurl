package config

import (
	"fmt"
	"time"
)

// Config holds all configuration options for murl
type Config struct {
	// Basic options
	Connections int           // Number of HTTP connections
	Duration    time.Duration // Duration of test
	Threads     int           // Number of threads
	Rate        int           // Requests per second (0 = unlimited)
	Timeout     time.Duration // Request timeout
	
	// Curl parsing
	CurlCommand string // Curl command to parse
	
	// HTTP options
	Method      string   // HTTP method
	Headers     []string // HTTP headers
	Body        string   // Request body
	ContentType string   // Content-Type header
	
	// Output options
	Verbose      bool // Verbose output
	PrintLatency bool // Print latency statistics
	
	// Engine options
	UseNetHTTP   bool // Force use standard library net/http instead of pulse
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.Connections <= 0 {
		return fmt.Errorf("connections must be greater than 0")
	}
	
	if c.Threads <= 0 {
		return fmt.Errorf("threads must be greater than 0")
	}
	
	if c.Duration <= 0 {
		return fmt.Errorf("duration must be greater than 0")
	}
	
	if c.Rate < 0 {
		return fmt.Errorf("rate cannot be negative")
	}
	
	if c.Timeout <= 0 {
		return fmt.Errorf("timeout must be greater than 0")
	}
	
	return nil
}

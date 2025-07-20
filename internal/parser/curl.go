package parser

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/antlabs/murl/internal/config"
	"github.com/antlabs/pcurl"
)

// ParseCurl parses a curl command and returns an http.Request
func ParseCurl(curlCommand string) (*http.Request, error) {
	return pcurl.ParseAndRequest(curlCommand)
}

// BuildRequest builds an http.Request from config and URL
func BuildRequest(cfg config.Config, targetURL *url.URL) (*http.Request, error) {
	var body strings.Reader
	if cfg.Body != "" {
		body = *strings.NewReader(cfg.Body)
	}
	
	req, err := http.NewRequest(cfg.Method, targetURL.String(), &body)
	if err != nil {
		return nil, err
	}
	
	// 添加headers
	for _, header := range cfg.Headers {
		parts := strings.SplitN(header, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			req.Header.Set(key, value)
		}
	}
	
	// 设置Content-Type
	if cfg.ContentType != "" {
		req.Header.Set("Content-Type", cfg.ContentType)
	} else if cfg.Body != "" && req.Header.Get("Content-Type") == "" {
		// 如果有body但没有设置Content-Type，默认设置为application/json
		req.Header.Set("Content-Type", "application/json")
	}
	
	// 设置User-Agent
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", "murl/1.0")
	}
	
	return req, nil
}

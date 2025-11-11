package api

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// HTTPServer wraps the HTTP server
type HTTPServer struct {
	server    *http.Server
	apiServer *Server
	port      int
}

// NewHTTPServer creates a new HTTP server
func NewHTTPServer(port int) *HTTPServer {
	apiServer := NewServer()

	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/api/v1/benchmark", apiServer.handleBenchmark)
	mux.HandleFunc("/api/v1/batch", apiServer.handleBatch)
	mux.HandleFunc("/api/v1/status/", apiServer.handleStatus)
	mux.HandleFunc("/api/v1/results/", apiServer.handleResults)

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok","service":"gurl-api"}`)
	})

	// Root endpoint with API info
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{
			"service": "gurl-api",
			"version": "1.0",
			"endpoints": {
				"benchmark": "POST /api/v1/benchmark",
				"batch": "POST /api/v1/batch",
				"status": "GET /api/v1/status/:id",
				"results": "GET /api/v1/results/:id",
				"health": "GET /health"
			}
		}`)
	})

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &HTTPServer{
		server:    server,
		apiServer: apiServer,
		port:      port,
	}
}

// Start starts the HTTP server
func (s *HTTPServer) Start() error {
	fmt.Printf("Starting gurl API server on port %d\n", s.port)
	fmt.Printf("API endpoints:\n")
	fmt.Printf("  POST   /api/v1/benchmark - Submit a benchmark task\n")
	fmt.Printf("  POST   /api/v1/batch     - Submit a batch test task\n")
	fmt.Printf("  GET    /api/v1/status/:id - Get task status\n")
	fmt.Printf("  GET    /api/v1/results/:id - Get task results\n")
	fmt.Printf("  GET    /health           - Health check\n")
	fmt.Printf("  GET    /                 - API information\n")

	return s.server.ListenAndServe()
}

// Stop stops the HTTP server gracefully
func (s *HTTPServer) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// StartWithContext starts the server and handles shutdown signals
func (s *HTTPServer) StartWithContext(ctx context.Context) error {
	// Start server in goroutine
	errChan := make(chan error, 1)
	go func() {
		if err := s.Start(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// Wait for context cancellation or error
	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		// Graceful shutdown
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return s.Stop(shutdownCtx)
	}
}

package mock

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

// ServerConfig holds mock server configuration
type ServerConfig struct {
	Port       int
	Delay      time.Duration
	Response   string
	StatusCode int
	Routes     []RouteConfig
	// EnableLogging controls whether per-request logs are printed
	EnableLogging bool
}

// RouteConfig defines a single route configuration
type RouteConfig struct {
	Path       string `json:"path" yaml:"path"`
	Method     string `json:"method" yaml:"method"`
	StatusCode int    `json:"status_code" yaml:"status_code"`
	Response   string `json:"response" yaml:"response"`
	Delay      string `json:"delay" yaml:"delay"`
	Echo       bool   `json:"echo" yaml:"echo"`
}

// Server represents a mock HTTP server
type Server struct {
	config ServerConfig
	server *http.Server
}

// NewServer creates a new mock server
func NewServer(config ServerConfig) *Server {
	return &Server{
		config: config,
	}
}

// Start starts the mock server
func (s *Server) Start() error {
	mux := http.NewServeMux()

	// 如果有配置文件中的路由，注册它们
	if len(s.config.Routes) > 0 {
		// 为了支持同一路径下的多种 HTTP 方法，这里先按 Path 归组，
		// 然后为每个 Path 注册一个 handler，在 handler 内部分发 Method，
		// 避免对同一 pattern 多次调用 HandleFunc 导致冲突。
		pathRoutes := make(map[string][]RouteConfig)
		for _, route := range s.config.Routes {
			path := route.Path
			if path == "" {
				path = "/"
			}
			pathRoutes[path] = append(pathRoutes[path], route)
		}

		for path, routes := range pathRoutes {
			s.registerRoute(mux, path, routes)
		}
	} else {
		// 默认路由：处理所有请求
		mux.HandleFunc("/", s.defaultHandler)
	}

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.Port),
		Handler: mux,
	}

	slog.Info("Mock server starting", "addr", fmt.Sprintf("http://localhost:%d", s.config.Port))
	slog.Info("Press Ctrl+C to stop")

	return s.server.ListenAndServe()
}

// registerRoute registers all routes for a specific path
func (s *Server) registerRoute(mux *http.ServeMux, path string, routes []RouteConfig) {
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Received request: %s %s\n", r.Method, r.URL.Path)
		for _, route := range routes {
			// 如果配置了 Method，则需要精确匹配；未配置则匹配所有方法
			if route.Method != "" && r.Method != route.Method {
				continue
			}
			handler := s.createRouteHandler(route)
			handler(w, r)
			return
		}

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	for _, route := range routes {
		if s.config.EnableLogging {
			slog.Info("Registered route", "method", route.Method, "path", path)
		}
	}
}

// createRouteHandler creates a handler for a specific route
func (s *Server) createRouteHandler(route RouteConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 解析延迟
		var delay time.Duration
		if route.Delay != "" {
			d, err := time.ParseDuration(route.Delay)
			if err == nil {
				delay = d
			}
		}

		// 应用延迟
		if delay > 0 {
			time.Sleep(delay)
		}

		// 设置状态码
		statusCode := route.StatusCode
		if statusCode == 0 {
			statusCode = 200
		}

		// 记录请求
		if s.config.EnableLogging {
			slog.Info("Request handled", "method", r.Method, "path", r.URL.Path, "status", statusCode, "delay", delay)
		}

		w.WriteHeader(statusCode)

		// Echo 模式：返回请求内容
		if route.Echo {
			s.echoRequest(w, r)
			return
		}

		// 自定义响应
		if route.Response != "" {
			w.Header().Set("Content-Type", "application/json")
			if _, err := fmt.Fprint(w, route.Response); err != nil {
				slog.Error("Error writing response", "err", err)
			}
			return
		}

		// 默认响应
		s.defaultResponse(w, r, statusCode)
	}
}

// defaultHandler handles all requests when no routes are configured
func (s *Server) defaultHandler(w http.ResponseWriter, r *http.Request) {
	// 应用延迟
	if s.config.Delay > 0 {
		time.Sleep(s.config.Delay)
	}

	// 记录请求
	if s.config.EnableLogging {
		slog.Info("Request handled", "method", r.Method, "path", r.URL.Path, "status", s.config.StatusCode, "delay", s.config.Delay)
	}

	w.WriteHeader(s.config.StatusCode)

	// 如果有自定义响应，使用它
	if s.config.Response != "" {
		w.Header().Set("Content-Type", "application/json")
		if _, err := fmt.Fprint(w, s.config.Response); err != nil {
			slog.Error("Error writing response", "err", err)
		}
		return
	}

	// 默认：echo 模式
	s.echoRequest(w, r)
}

// echoRequest echoes the request back to the client
func (s *Server) echoRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 读取请求体
	body, _ := io.ReadAll(r.Body)
	defer func() {
		if err := r.Body.Close(); err != nil {
			slog.Error("Error closing request body", "err", err)
		}
	}()

	// 构建响应
	response := map[string]interface{}{
		"method":  r.Method,
		"path":    r.URL.Path,
		"query":   r.URL.RawQuery,
		"headers": r.Header,
		"body":    string(body),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("Error encoding response", "err", err)
	}
}

// defaultResponse provides a default JSON response
func (s *Server) defaultResponse(w http.ResponseWriter, r *http.Request, statusCode int) {
	w.Header().Set("Content-Type", "application/json")

	response := map[string]interface{}{
		"status":  statusCode,
		"message": http.StatusText(statusCode),
		"path":    r.URL.Path,
		"method":  r.Method,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("Error encoding response", "err", err)
	}
}

// Stop stops the mock server
func (s *Server) Stop() error {
	if s.server != nil {
		return s.server.Close()
	}
	return nil
}

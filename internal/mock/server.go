package mock

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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
		for _, route := range s.config.Routes {
			s.registerRoute(mux, route)
		}
	} else {
		// 默认路由：处理所有请求
		mux.HandleFunc("/", s.defaultHandler)
	}

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.Port),
		Handler: mux,
	}

	log.Printf("Mock server starting on http://localhost:%d\n", s.config.Port)
	log.Printf("Press Ctrl+C to stop\n")

	return s.server.ListenAndServe()
}

// registerRoute registers a single route
func (s *Server) registerRoute(mux *http.ServeMux, route RouteConfig) {
	handler := s.createRouteHandler(route)
	
	if route.Method != "" {
		// 如果指定了方法，只处理该方法
		mux.HandleFunc(route.Path, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != route.Method {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}
			handler(w, r)
		})
	} else {
		mux.HandleFunc(route.Path, handler)
	}

	log.Printf("Registered route: %s %s", route.Method, route.Path)
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
		log.Printf("%s %s - Status: %d, Delay: %v", r.Method, r.URL.Path, statusCode, delay)

		w.WriteHeader(statusCode)

		// Echo 模式：返回请求内容
		if route.Echo {
			s.echoRequest(w, r)
			return
		}

		// 自定义响应
		if route.Response != "" {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, route.Response)
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
	log.Printf("%s %s - Status: %d, Delay: %v", r.Method, r.URL.Path, s.config.StatusCode, s.config.Delay)

	w.WriteHeader(s.config.StatusCode)

	// 如果有自定义响应，使用它
	if s.config.Response != "" {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, s.config.Response)
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
	defer r.Body.Close()

	// 构建响应
	response := map[string]interface{}{
		"method":  r.Method,
		"path":    r.URL.Path,
		"query":   r.URL.RawQuery,
		"headers": r.Header,
		"body":    string(body),
	}

	json.NewEncoder(w).Encode(response)
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

	json.NewEncoder(w).Encode(response)
}

// Stop stops the mock server
func (s *Server) Stop() error {
	if s.server != nil {
		return s.server.Close()
	}
	return nil
}

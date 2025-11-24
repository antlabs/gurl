package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/antlabs/gurl/internal/api"
	"github.com/antlabs/gurl/internal/batch"
	"github.com/antlabs/gurl/internal/benchmark"
	"github.com/antlabs/gurl/internal/compare"
	"github.com/antlabs/gurl/internal/config"
	"github.com/antlabs/gurl/internal/mcp"
	"github.com/antlabs/gurl/internal/mock"
	"github.com/antlabs/gurl/internal/parser"
	"github.com/antlabs/gurl/internal/template"
	"github.com/guonaihong/clop"
)

// Args 定义命令行参数结构
type Args struct {
	// 基本选项
	Connections int           `clop:"-c;--connections" usage:"Number of HTTP connections to keep open" default:"10"`
	Duration    time.Duration `clop:"-d;--duration" usage:"Duration of test" default:"10s"`
	Threads     int           `clop:"-t;--threads" usage:"Number of threads to use" default:"2"`
	Rate        int           `clop:"-R;--rate" usage:"Work rate (requests/sec) 0=unlimited" default:"0"`
	Timeout     time.Duration `clop:"--timeout" usage:"Socket/request timeout" default:"30s"`
	Requests    int64         `clop:"-n;--requests" usage:"Total number of requests to perform (0=unlimited, duration-limited)" default:"0"`

	// curl解析选项
	CurlCommand  string `clop:"--parse-curl" usage:"Parse curl command and use it for benchmarking"`
	CurlFile     string `clop:"--parse-curl-file" usage:"Parse multiple curl commands from file (one per line)"`
	LoadStrategy string `clop:"--load-strategy" usage:"Load distribution strategy: random, round-robin" default:"random"`

	// HTTP选项
	Method      string   `clop:"-X;--method" usage:"HTTP method" default:"GET"`
	Headers     []string `clop:"-H;--header" usage:"HTTP header to add to request"`
	Body        string   `clop:"--data" usage:"HTTP request body"`
	ContentType string   `clop:"--content-type" usage:"Content-Type header"`

	// 输出选项
	Verbose      bool   `clop:"-v;--verbose" usage:"Verbose output"`
	PrintLatency bool   `clop:"--latency" usage:"Print latency statistics"`
	LiveUI       bool   `clop:"--live-ui" usage:"Enable live terminal UI with real-time stats"`
	UITheme      string `clop:"--ui-theme" usage:"UI color theme: dark, light, or auto (default: auto)"`

	// 引擎选项
	UseNetHTTP bool `clop:"--use-nethttp" usage:"Force use standard library net/http instead of pulse"`

	// 批量测试选项
	BatchConfig      string `clop:"--batch-config" usage:"Path to batch test configuration file (YAML/JSON)"`
	BatchConcurrency int    `clop:"--batch-concurrency" usage:"Maximum concurrent batch tests" default:"3"`
	BatchSequential  bool   `clop:"--batch-sequential" usage:"Run batch tests sequentially instead of concurrently"`
	BatchReport      string `clop:"--batch-report" usage:"Output format for batch report (text|csv|json)" default:"text"`

	// 模板变量选项
	Variables     []string `clop:"--var" usage:"Define template variables (format: name=type:params)"`
	HelpTemplates bool     `clop:"--help-templates" usage:"Show template variable help and examples"`

	// MCP选项
	MCP         bool   `clop:"--mcp" usage:"Start as an MCP server"`
	MCPDebugLog string `clop:"--mcp-debug-log" usage:"Path to MCP debug log file (only used with --mcp)"`

	// API服务器选项
	APIServer bool `clop:"--api-server" usage:"Start as a RESTful API server"`
	APIPort   int  `clop:"--api-port" usage:"Port for API server" default:"8080"`

	// Mock服务器选项
	MockServer     bool   `clop:"--mock-server" usage:"Start a mock HTTP server for testing"`
	MockPort       int    `clop:"--mock-port" usage:"Port for mock server" default:"8080"`
	MockDelay      string `clop:"--mock-delay" usage:"Response delay (e.g., 100ms, 1s)" default:"0s"`
	MockResponse   string `clop:"--mock-response" usage:"Custom response body"`
	MockStatusCode int    `clop:"--mock-status" usage:"HTTP status code to return" default:"200"`
	MockConfig     string `clop:"--mock-config" usage:"Path to mock server configuration file (YAML/JSON)"`

	// Compare 选项（阶段1：使用 --compare-config/-f 与 --compare-name/-n）
	CompareConfig string `clop:"--compare-config;-f" usage:"Path to compare configuration file (YAML/JSON)"`
	CompareName   string `clop:"--compare-name" usage:"Name of compare scenario to run"`

	// 位置参数
	URL string `clop:"args=url" usage:"Target URL for benchmarking"`
}

// argsToConfig 将Args转换为config.Config
func (a *Args) toConfig() config.Config {
	return config.Config{
		Connections:  a.Connections,
		Duration:     a.Duration,
		Threads:      a.Threads,
		Rate:         a.Rate,
		Timeout:      a.Timeout,
		Requests:     a.Requests,
		CurlCommand:  a.CurlCommand,
		CurlFile:     a.CurlFile,
		LoadStrategy: a.LoadStrategy,
		Method:       a.Method,
		Headers:      a.Headers,
		Body:         a.Body,
		ContentType:  a.ContentType,
		Verbose:      a.Verbose,
		PrintLatency: a.PrintLatency,
		LiveUI:       a.LiveUI,
		UITheme:      a.UITheme,
		UseNetHTTP:   a.UseNetHTTP,
	}
}

// runBenchmark 执行基准测试
func runBenchmark(args *Args) error {
	var req *http.Request
	var requests []*http.Request
	var err error

	cfg := args.toConfig()

	// 创建模板解析器并设置变量
	templateParser := template.NewTemplateParser()
	if len(args.Variables) > 0 {
		context, err := template.ParseVariableDefinitions(args.Variables)
		if err != nil {
			return fmt.Errorf("failed to parse variable definitions: %w", err)
		}
		templateParser = template.NewTemplateParserWithContext(context)
	}

	// 处理多个curl命令文件
	if args.CurlFile != "" {
		requests, err = parser.ParseCurlFile(args.CurlFile)
		if err != nil {
			return fmt.Errorf("failed to parse curl file: %w", err)
		}
		fmt.Printf("Loaded %d curl commands from file\n", len(requests))
		fmt.Printf("Load strategy: %s\n", cfg.LoadStrategy)
	} else if args.CurlCommand != "" {
		// 处理模板变量
		processedCurl := args.CurlCommand
		if len(args.Variables) > 0 || template.HasTemplateVariables(args.CurlCommand) {
			processedCurl, err = templateParser.ParseTemplate(args.CurlCommand)
			if err != nil {
				return fmt.Errorf("failed to process template variables in curl command: %w", err)
			}
			if args.Verbose {
				fmt.Printf("Original curl: %s\n", args.CurlCommand)
				fmt.Printf("Processed curl: %s\n", processedCurl)
			}
		}

		// 解析curl命令
		req, err = parser.ParseCurl(processedCurl)
		if err != nil {
			return fmt.Errorf("failed to parse curl command: %w", err)
		}
	} else {
		// 使用传统方式构建请求
		if args.URL == "" {
			return fmt.Errorf("URL is required when not using --parse-curl")
		}

		// 处理URL中的模板变量
		targetURL := args.URL
		if len(args.Variables) > 0 || template.HasTemplateVariables(args.URL) {
			targetURL, err = templateParser.ParseTemplate(args.URL)
			if err != nil {
				return fmt.Errorf("failed to process template variables in URL: %w", err)
			}
			if args.Verbose {
				fmt.Printf("Original URL: %s\n", args.URL)
				fmt.Printf("Processed URL: %s\n", targetURL)
			}
		}

		if !strings.HasPrefix(targetURL, "http://") && !strings.HasPrefix(targetURL, "https://") {
			targetURL = "http://" + targetURL
		}

		parsedURL, err := url.Parse(targetURL)
		if err != nil {
			return fmt.Errorf("invalid URL: %w", err)
		}

		req, err = parser.BuildRequest(cfg, parsedURL)
		if err != nil {
			return fmt.Errorf("failed to build request: %w", err)
		}
	}

	// 验证配置
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// 创建上下文用于优雅关闭
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 处理信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nShutting down...")
		cancel()
	}()

	// 创建并运行基准测试
	var bench *benchmark.Benchmark
	var targetURL string

	if len(requests) > 0 {
		// 多请求模式
		bench = benchmark.NewWithMultipleRequests(cfg, requests)
		targetURL = fmt.Sprintf("%d endpoints", len(requests))
	} else {
		// 单请求模式
		bench = benchmark.New(cfg, req)
		targetURL = req.URL.String()
	}

	// 只在非 LiveUI 模式下输出初始信息
	if !cfg.LiveUI {
		fmt.Printf("Running %s test @ %s\n", cfg.Duration, targetURL)
		fmt.Printf("  %d threads and %d connections\n", cfg.Threads, cfg.Connections)
	}

	results, err := bench.Run(ctx)
	if err != nil {
		return fmt.Errorf("benchmark failed: %w", err)
	}

	// 打印结果
	benchmark.PrintResults(results, cfg)

	return nil
}

// runCompare 执行 compare 模式（一对一场景）
func runCompare(args *Args) error {
	if args.CompareConfig == "" {
		return fmt.Errorf("compare-config is required")
	}
	if args.CompareName == "" {
		return fmt.Errorf("compare-name is required")
	}

	// 加载 compare 配置
	cmpCfg, err := config.LoadCompareConfig(args.CompareConfig)
	if err != nil {
		return fmt.Errorf("failed to load compare config: %w", err)
	}

	// 执行指定场景
	results, passed, failed, err := compare.RunScenario(cmpCfg, args.CompareName)
	if err != nil {
		return fmt.Errorf("failed to run compare scenario: %w", err)
	}

	// 查找场景以获得 base/target 名称
	var scenario *config.CompareScenario
	for i := range cmpCfg.Scenarios {
		if cmpCfg.Scenarios[i].Name == args.CompareName {
			scenario = &cmpCfg.Scenarios[i]
			break
		}
	}

	fmt.Printf("Scenario: %s\n", args.CompareName)
	if scenario != nil {
		fmt.Printf("Mode   : %s\n", scenario.Mode)
		if scenario.Base != "" {
			fmt.Printf("Base   : %s\n", scenario.Base)
		}
		if scenario.Target != "" {
			fmt.Printf("Target : %s\n", scenario.Target)
		}
	}
	fmt.Println()

	// 按 PairLabel 对结果分组，便于批量模式下分块输出
	pairs := make(map[string][]compare.AssertionResult)
	order := make([]string, 0)
	for _, r := range results {
		label := r.PairLabel
		if label == "" {
			label = "(single pair)"
		}
		if _, ok := pairs[label]; !ok {
			order = append(order, label)
		}
		pairs[label] = append(pairs[label], r)
	}

	for _, label := range order {
		fmt.Printf("=== Pair: %s ===\n", label)
		for _, r := range pairs[label] {
			status := "FAIL"
			if r.OK {
				status = "OK"
			}
			fmt.Printf("[%s] %s\n", status, r.Line)
			if r.BaseValue != "" || r.TargetValue != "" {
				if r.BaseValue != "" {
					fmt.Printf("     base:   %s\n", r.BaseValue)
				}
				if r.TargetValue != "" {
					fmt.Printf("     target: %s\n", r.TargetValue)
				}
			}
			if !r.OK && r.Message != "" {
				fmt.Printf("     reason: %s\n", r.Message)
			}
		}
		fmt.Println()
	}

	fmt.Printf("Summary: %d passed, %d failed\n", passed, failed)
	if failed > 0 {
		return fmt.Errorf("compare scenario failed")
	}
	return nil
}

// runBatchTest 执行批量测试
func runBatchTest(args *Args) error {
	// 加载批量配置文件
	batchConfig, err := config.LoadBatchConfig(args.BatchConfig)
	if err != nil {
		return fmt.Errorf("failed to load batch config: %w", err)
	}

	// 创建默认配置
	defaults := args.toConfig()

	// 创建批量执行器
	executor := batch.NewExecutor(args.BatchConcurrency, args.Verbose)

	// 创建上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 处理信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nReceived interrupt signal, stopping batch tests...")
		cancel()
	}()

	// 执行批量测试
	var result *batch.BatchResult
	if args.BatchSequential {
		result, err = executor.ExecuteSequential(ctx, batchConfig, &defaults)
	} else {
		result, err = executor.Execute(ctx, batchConfig, &defaults)
	}

	if err != nil {
		return fmt.Errorf("batch test failed: %w", err)
	}

	// 生成报告
	reporter := batch.NewReporter(args.Verbose)

	// 根据输出格式生成报告
	switch args.BatchReport {
	case "csv":
		fmt.Print(reporter.GenerateCSVReport(result))
	case "json":
		fmt.Print(reporter.GenerateJSONReport(result))
	default: // "text"
		fmt.Print(reporter.GenerateReport(result))
	}

	// 打印简要摘要
	reporter.PrintSummary(result)

	return nil
}

// runMockServer 启动 mock HTTP 服务器
func runMockServer(args *Args) error {
	var serverConfig mock.ServerConfig

	// 如果有配置文件，加载它
	if args.MockConfig != "" {
		config, err := mock.LoadConfig(args.MockConfig)
		if err != nil {
			return fmt.Errorf("failed to load mock config: %w", err)
		}

		serverConfig.Port = config.Port
		if serverConfig.Port == 0 {
			serverConfig.Port = args.MockPort
		}
		serverConfig.Routes = config.Routes
	} else {
		// 使用命令行参数
		serverConfig.Port = args.MockPort
		serverConfig.StatusCode = args.MockStatusCode
		serverConfig.Response = args.MockResponse

		// 解析延迟
		if args.MockDelay != "" {
			delay, err := time.ParseDuration(args.MockDelay)
			if err != nil {
				return fmt.Errorf("invalid delay format: %w", err)
			}
			serverConfig.Delay = delay
		}
	}

	// 创建并启动服务器
	server := mock.NewServer(serverConfig)

	// 处理信号以优雅关闭
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nShutting down mock server...")
		if err := server.Stop(); err != nil {
			fmt.Fprintf(os.Stderr, "Error stopping server: %v\n", err)
		}
		os.Exit(0)
	}()

	return server.Start()
}

// Execute 执行命令行程序
func main() {
	args := &Args{}

	if err := clop.Bind(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing arguments: %v\n", err)
		os.Exit(1)
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError})))

	// 检查是否启动 Mock 服务器
	if args.MockServer {
		if err := runMockServer(args); err != nil {
			fmt.Fprintf(os.Stderr, "Mock server error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// 检查是否启动MCP服务
	if args.MCP {
		server := mcp.NewServer()
		if err := server.StartWithDebugLog(args.MCPDebugLog); err != nil {
			fmt.Fprintf(os.Stderr, "MCP server error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// 检查是否启动API服务器
	if args.APIServer {
		apiServer := api.NewHTTPServer(args.APIPort)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// 处理信号以优雅关闭
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-sigChan
			fmt.Println("\nShutting down API server...")
			cancel()
		}()

		if err := apiServer.StartWithContext(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "API server error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// 检查是否显示模板帮助
	if args.HelpTemplates {
		fmt.Print(template.GetAllFunctionsHelp())
		fmt.Print("\n\n")
		fmt.Print(template.PrintTemplateExamples())
		fmt.Print("\n\n")
		fmt.Print(template.GetQuickStartGuide())
		return
	}

	// 检查是否为 compare 模式（优先于批量/基准测试）
	if args.CompareConfig != "" {
		if err := runCompare(args); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// 检查是否为批量测试模式
	if args.BatchConfig != "" {
		if err := runBatchTest(args); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	} else {
		// 单个测试模式
		if err := runBenchmark(args); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}
}

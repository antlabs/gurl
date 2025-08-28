package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/antlabs/gurl/internal/batch"
	"github.com/antlabs/gurl/internal/benchmark"
	"github.com/antlabs/gurl/internal/config"
	"github.com/antlabs/gurl/internal/mcp"
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

	// curl解析选项
	CurlCommand string `clop:"--parse-curl" usage:"Parse curl command and use it for benchmarking"`

	// HTTP选项
	Method      string   `clop:"-X;--method" usage:"HTTP method" default:"GET"`
	Headers     []string `clop:"-H;--header" usage:"HTTP header to add to request"`
	Body        string   `clop:"--data" usage:"HTTP request body"`
	ContentType string   `clop:"--content-type" usage:"Content-Type header"`

	// 输出选项
	Verbose      bool `clop:"-v;--verbose" usage:"Verbose output"`
	PrintLatency bool `clop:"--latency" usage:"Print latency statistics"`

	// 引擎选项
	UseNetHTTP bool `clop:"--use-nethttp" usage:"Force use standard library net/http instead of pulse"`

	// 批量测试选项
	BatchConfig     string `clop:"--batch-config" usage:"Path to batch test configuration file (YAML/JSON)"`
	BatchConcurrency int   `clop:"--batch-concurrency" usage:"Maximum concurrent batch tests" default:"3"`
	BatchSequential bool   `clop:"--batch-sequential" usage:"Run batch tests sequentially instead of concurrently"`
	BatchReport     string `clop:"--batch-report" usage:"Output format for batch report (text|csv|json)" default:"text"`

	// 模板变量选项
	Variables    []string `clop:"--var" usage:"Define template variables (format: name=type:params)"`
	HelpTemplates bool     `clop:"--help-templates" usage:"Show template variable help and examples"`

	// MCP选项
	MCP bool `clop:"--mcp" usage:"Start as an MCP server"`

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
		CurlCommand:  a.CurlCommand,
		Method:       a.Method,
		Headers:      a.Headers,
		Body:         a.Body,
		ContentType:  a.ContentType,
		Verbose:      a.Verbose,
		PrintLatency: a.PrintLatency,
		UseNetHTTP:   a.UseNetHTTP,
	}
}

// runBenchmark 执行基准测试
func runBenchmark(args *Args) error {
	var req *http.Request
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

	// 处理URL和curl命令解析
	if args.CurlCommand != "" {
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

	// 创建并运行基准测试（自动选择pulse或net/http方式）
	bench := benchmark.New(cfg, req)

	fmt.Printf("Running %s test @ %s\n", cfg.Duration, req.URL.String())
	fmt.Printf("  %d threads and %d connections\n", cfg.Threads, cfg.Connections)

	results, err := bench.Run(ctx)
	if err != nil {
		return fmt.Errorf("benchmark failed: %w", err)
	}

	// 打印结果
	benchmark.PrintResults(results, cfg)

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

// Execute 执行命令行程序
func main() {
	args := &Args{}

	clop.Bind(args)

	// 检查是否启动MCP服务
	if args.MCP {
		server := mcp.NewServer()
		if err := server.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "MCP server error: %v\n", err)
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

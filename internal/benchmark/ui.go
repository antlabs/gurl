package benchmark

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// ColorTheme defines color scheme for terminal UI
type ColorTheme struct {
	Name string

	// Border and text colors
	Border      ui.Color
	Text        ui.Color
	Title       ui.Color
	ProgressBar ui.Color

	// Stats colors
	RequestsLabel ui.Color
	Slowest       ui.Color
	Fastest       ui.Color
	Average       ui.Color
	ReqSec        ui.Color

	// Status code colors
	Status2xx ui.Color
	Status3xx ui.Color
	Status4xx ui.Color
	Status5xx ui.Color
	Error     ui.Color

	// Chart colors
	ReqChartBar    ui.Color
	ReqChartLabel  ui.Color
	ReqChartNumber ui.Color

	LatencyChartBar    ui.Color
	LatencyChartLabel  ui.Color
	LatencyChartNumber ui.Color

	// Help text color
	Help ui.Color
}

// LiveUI represents a live terminal UI for benchmark progress
type LiveUI struct {
	mu sync.RWMutex

	// Widgets
	progressGauge *widgets.Gauge
	statsTable    *widgets.Paragraph
	statusTable   *widgets.Paragraph
	reqChart      *widgets.BarChart
	latencyChart  *widgets.BarChart
	endpointTable *widgets.Table // 新增：端点统计表格
	helpText      *widgets.Paragraph

	// Data
	totalRequests int64
	reqPerSecond  []int64
	statusCodes   map[int]int64
	endpointStats map[string]*EndpointLiveStats // 新增：端点实时统计

	startTime     time.Time
	duration      time.Duration
	multiEndpoint bool // 是否为多端点模式

	// Theme
	theme *ColorTheme

	// Control
	stopChan chan struct{}
}

// EndpointLiveStats holds live statistics for a single endpoint
type EndpointLiveStats struct {
	URL        string
	Requests   int64
	ReqPerSec  float64
	AvgLatency time.Duration
	MinLatency time.Duration
	MaxLatency time.Duration
	Errors     int64
	LastUpdate time.Time
}

// DarkTheme returns a color theme optimized for dark terminal backgrounds
func DarkTheme() *ColorTheme {
	return &ColorTheme{
		Name:        "dark",
		Border:      ui.ColorWhite,
		Text:        ui.ColorWhite,
		Title:       ui.ColorWhite,
		ProgressBar: ui.ColorBlue,

		RequestsLabel: ui.ColorWhite,
		Slowest:       ui.ColorRed,
		Fastest:       ui.ColorGreen,
		Average:       ui.ColorCyan,
		ReqSec:        ui.ColorYellow,

		Status2xx: ui.ColorGreen,
		Status3xx: ui.ColorCyan,
		Status4xx: ui.ColorYellow,
		Status5xx: ui.ColorRed,
		Error:     ui.ColorRed,

		ReqChartBar:    ui.ColorGreen,
		ReqChartLabel:  ui.ColorWhite,
		ReqChartNumber: ui.ColorYellow,

		LatencyChartBar:    ui.ColorYellow,
		LatencyChartLabel:  ui.ColorWhite,
		LatencyChartNumber: ui.ColorCyan,

		Help: ui.ColorYellow,
	}
}

// LightTheme returns a color theme optimized for light terminal backgrounds
// Paper主题配色：基于 #eeeeee 背景的专业配色方案
func LightTheme() *ColorTheme {
	// Paper 主题配色
	// 背景 255 (#eeeeee)
	// 一级文字 238 (#444444) - 深灰
	// 二级文字 244 (#808080) - 中灰
	// 高亮 166 (#df5f00) - 橙色
	// 成功 64 (#008700) - 绿色
	// 警告 172 (#d75f00) - 深橙
	// 危险 160 (#d70000) - 红色
	// 链接 25 (#005fd7) - 蓝色

	return &ColorTheme{
		Name:        "light",
		Border:      25,  // 蓝色边框 (#005fd7) - 清晰界限
		Text:        238, // 深灰文本 (#444444) - 主要内容
		Title:       25,  // 蓝色标题 (#005fd7) - 醒目标识
		ProgressBar: 166, // 橙色进度条 (#df5f00) - 高亮进度

		RequestsLabel: 238, // 深灰标签 (#444444) - 主要文字
		Slowest:       160, // 红色最慢 (#d70000) - 危险警告
		Fastest:       64,  // 绿色最快 (#008700) - 成功标识
		Average:       166, // 橙色平均 (#df5f00) - 高亮关键
		ReqSec:        25,  // 蓝色速率 (#005fd7) - 重要指标

		Status2xx: 64,  // 绿色成功 (#008700) - 正常状态
		Status3xx: 25,  // 蓝色重定向 (#005fd7) - 中性状态
		Status4xx: 172, // 深橙客户端错误 (#d75f00) - 警告状态
		Status5xx: 160, // 红色服务器错误 (#d70000) - 危险状态
		Error:     160, // 红色错误 (#d70000) - 严重问题

		ReqChartBar:    64,  // 绿色柱状图 (#008700) - 成功表现
		ReqChartLabel:  238, // 深灰标签 (#444444) - 清晰可读
		ReqChartNumber: 238, // 深灰数字 (#444444) - 清晰可读

		LatencyChartBar:    166, // 橙色柱状图 (#df5f00) - 高亮延迟
		LatencyChartLabel:  238, // 深灰标签 (#444444) - 清晰可读
		LatencyChartNumber: 238, // 深灰数字 (#444444) - 清晰可读

		Help: 172, // 深橙帮助文本 (#d75f00) - 警示引导
	}
}

// DetectTerminalTheme tries to detect whether the terminal has a dark or light background
// Returns "dark" or "light"
func DetectTerminalTheme() string {
	// Priority 1: Check for explicit theme environment variable (custom)
	if theme := os.Getenv("GURL_THEME"); theme != "" {
		theme = strings.ToLower(theme)
		if theme == "light" || theme == "dark" {
			return theme
		}
	}

	// Priority 2: Check COLORFGBG environment variable
	// COLORFGBG format: "foreground;background" where 0-6 is dark, 7-15 is light
	if colorfgbg := os.Getenv("COLORFGBG"); colorfgbg != "" {
		parts := strings.Split(colorfgbg, ";")
		if len(parts) >= 2 {
			bg := parts[len(parts)-1]
			// Background color codes: 0-6 are dark, 7-15 are light
			if bg == "0" || bg == "8" {
				return "dark"
			} else if bg == "7" || bg == "15" {
				return "light"
			}
		}
	}

	// Priority 3: Detect based on terminal program and known defaults
	termProgram := os.Getenv("TERM_PROGRAM")
	term := os.Getenv("TERM")

	// macOS Terminal.app - default is light theme
	if termProgram == "Apple_Terminal" {
		// Try to detect via terminal profile name (if available)
		// This is a heuristic - Terminal.app default profiles are usually light
		return "light"
	}

	// VSCode terminal - check COLOR_THEME if available
	if colorTheme := os.Getenv("COLOR_THEME"); colorTheme != "" {
		colorTheme = strings.ToLower(colorTheme)
		if strings.Contains(colorTheme, "light") || strings.Contains(colorTheme, "default") {
			return "light"
		}
		if strings.Contains(colorTheme, "dark") {
			return "dark"
		}
	}

	// Windows Terminal - check WT_SESSION
	if os.Getenv("WT_SESSION") != "" {
		// Windows Terminal default is usually dark, but can be configured
		// Check WT_PROFILE_ID for profile name hints
		if profileID := os.Getenv("WT_PROFILE_ID"); profileID != "" {
			profileID = strings.ToLower(profileID)
			if strings.Contains(profileID, "light") {
				return "light"
			}
		}
		// Default to dark for Windows Terminal
		return "dark"
	}

	// iTerm2 - check ITERM_PROFILE or ITERM_SESSION_ID
	if strings.Contains(termProgram, "iTerm") || os.Getenv("ITERM_SESSION_ID") != "" {
		// iTerm2 default is usually dark, but can be configured
		// Check profile name if available
		if profile := os.Getenv("ITERM_PROFILE"); profile != "" {
			profile = strings.ToLower(profile)
			if strings.Contains(profile, "light") || strings.Contains(profile, "solarized light") {
				return "light"
			}
		}
		// Default to dark for iTerm2
		return "dark"
	}

	// GNOME Terminal - check for common light theme indicators
	if term != "" && (strings.Contains(term, "gnome") || strings.Contains(term, "xterm")) {
		// Check for common environment variables that might indicate theme
		// This is heuristic-based
		if desktopSession := os.Getenv("DESKTOP_SESSION"); desktopSession != "" {
			desktopSession = strings.ToLower(desktopSession)
			// Some desktop environments default to light themes
			if strings.Contains(desktopSession, "gnome") || strings.Contains(desktopSession, "ubuntu") {
				// Could be either, but many default to dark now
			}
		}
	}

	// Priority 4: Default to dark theme (most developers use dark terminals)
	return "dark"
}

// GetTheme returns the appropriate theme based on detection or explicit setting
func GetTheme(explicitTheme string) *ColorTheme {
	theme := strings.ToLower(strings.TrimSpace(explicitTheme))
	// Treat "auto" and empty string as auto-detect
	if theme == "" || theme == "auto" {
		theme = DetectTerminalTheme()
	}

	if theme == "light" {
		return LightTheme()
	}
	return DarkTheme()
}

// NewLiveUI creates a new live UI with automatic theme detection
func NewLiveUI(duration time.Duration) (*LiveUI, error) {
	return NewLiveUIWithTheme(duration, "")
}

// NewLiveUIWithTheme creates a new live UI with explicit theme ("dark" or "light")
func NewLiveUIWithTheme(duration time.Duration, themeName string) (*LiveUI, error) {
	if err := ui.Init(); err != nil {
		return nil, err
	}

	// Get theme
	theme := GetTheme(themeName)

	liveUI := &LiveUI{
		statusCodes:   make(map[int]int64),
		endpointStats: make(map[string]*EndpointLiveStats),
		startTime:     time.Now(),
		duration:      duration,
		theme:         theme,
		stopChan:      make(chan struct{}),
	}

	// Progress gauge
	liveUI.progressGauge = widgets.NewGauge()
	liveUI.progressGauge.Title = "Progress"
	liveUI.progressGauge.SetRect(0, 0, 100, 3)
	liveUI.progressGauge.BarColor = theme.ProgressBar
	liveUI.progressGauge.BorderStyle.Fg = theme.Border

	// Stats table
	liveUI.statsTable = widgets.NewParagraph()
	liveUI.statsTable.Title = "Stats for last sec"
	liveUI.statsTable.SetRect(0, 3, 50, 12)
	liveUI.statsTable.BorderStyle.Fg = theme.Border

	// Status code table
	liveUI.statusTable = widgets.NewParagraph()
	liveUI.statusTable.Title = "Status code distribution"
	liveUI.statusTable.SetRect(50, 3, 100, 12)
	liveUI.statusTable.BorderStyle.Fg = theme.Border

	// Request chart
	liveUI.reqChart = widgets.NewBarChart()
	liveUI.reqChart.Title = "Requests / past sec (auto)"
	liveUI.reqChart.SetRect(0, 12, 50, 30)
	liveUI.reqChart.BarWidth = 7 // 增加宽度以显示更多数字
	liveUI.reqChart.BarGap = 1   // 柱子间隔
	liveUI.reqChart.BarColors = []ui.Color{theme.ReqChartBar}
	liveUI.reqChart.LabelStyles = []ui.Style{ui.NewStyle(theme.ReqChartLabel)}
	liveUI.reqChart.NumStyles = []ui.Style{ui.NewStyle(theme.ReqChartNumber)}
	liveUI.reqChart.BorderStyle.Fg = theme.Border

	// Latency chart
	liveUI.latencyChart = widgets.NewBarChart()
	liveUI.latencyChart.Title = "Response time histogram (ms)"
	liveUI.latencyChart.SetRect(50, 12, 100, 30)
	liveUI.latencyChart.BarWidth = 7 // 增加宽度
	liveUI.latencyChart.BarGap = 1   // 柱子间隔
	liveUI.latencyChart.BarColors = []ui.Color{theme.LatencyChartBar}
	liveUI.latencyChart.LabelStyles = []ui.Style{ui.NewStyle(theme.LatencyChartLabel)}
	liveUI.latencyChart.NumStyles = []ui.Style{ui.NewStyle(theme.LatencyChartNumber)}
	liveUI.latencyChart.BorderStyle.Fg = theme.Border

	// Endpoint table (initially hidden, shown when multi-endpoint mode)
	liveUI.endpointTable = widgets.NewTable()
	liveUI.endpointTable.Title = "Per-Endpoint Statistics (live)"
	liveUI.endpointTable.SetRect(0, 30, 100, 42)
	liveUI.endpointTable.TextStyle = ui.NewStyle(theme.Text)
	liveUI.endpointTable.BorderStyle.Fg = theme.Border
	liveUI.endpointTable.RowSeparator = false
	liveUI.endpointTable.FillRow = true

	// Help text
	liveUI.helpText = widgets.NewParagraph()
	helpColorName := "yellow"
	if theme.Name == "light" {
		// Paper主题使用深橙色 (172) 作为帮助文本
		helpColorName = "172"
	}
	liveUI.helpText.Text = fmt.Sprintf("[Press 'q' or 'Ctrl+C' to stop test early](fg:%s,mod:bold)", helpColorName)
	liveUI.helpText.SetRect(0, 42, 100, 45)
	liveUI.helpText.Border = false
	liveUI.helpText.TextStyle.Fg = theme.Help

	// 启动键盘事件监听
	go liveUI.handleKeyEvents()

	// 初始渲染一次，避免黑屏
	liveUI.Render()

	return liveUI, nil
}

// Update updates the UI with new data
func (l *LiveUI) Update(requests int64, reqPerSec int64, statusCodes map[int]int64, avgLatency, minLatency, maxLatency time.Duration, latencyPercentiles map[float64]time.Duration, errors int64) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.totalRequests = requests
	l.reqPerSecond = append(l.reqPerSecond, reqPerSec)
	l.statusCodes = statusCodes

	// Keep only last 10 samples for chart
	if len(l.reqPerSecond) > 10 {
		l.reqPerSecond = l.reqPerSecond[len(l.reqPerSecond)-10:]
	}

	// Update progress
	elapsed := time.Since(l.startTime)
	progress := int(float64(elapsed) / float64(l.duration) * 100)
	if progress > 100 {
		progress = 100
	}
	l.progressGauge.Percent = progress
	l.progressGauge.Label = fmt.Sprintf("%ds / %ds", int(elapsed.Seconds()), int(l.duration.Seconds()))

	// Update stats with theme colors
	requestsColor := l.getColorName(l.theme.RequestsLabel)
	slowestColor := l.getColorName(l.theme.Slowest)
	fastestColor := l.getColorName(l.theme.Fastest)
	averageColor := l.getColorName(l.theme.Average)
	reqSecColor := l.getColorName(l.theme.ReqSec)

	l.statsTable.Text = fmt.Sprintf(
		"[Requests](fg:%s,mod:bold) : %d\n"+
			"[Slowest](fg:%s) : %.4f secs\n"+
			"[Fastest](fg:%s) : %.4f secs\n"+
			"[Average](fg:%s) : %.4f secs\n"+
			"[Req/Sec](fg:%s) : %.2f",
		requestsColor, requests,
		slowestColor, maxLatency.Seconds(),
		fastestColor, minLatency.Seconds(),
		averageColor, avgLatency.Seconds(),
		reqSecColor, float64(reqPerSec),
	)

	// Update status codes with theme colors
	statusText := ""
	if len(statusCodes) > 0 {
		for code, count := range statusCodes {
			// 根据状态码类型设置颜色
			var color string
			if code >= 200 && code < 300 {
				color = l.getColorName(l.theme.Status2xx)
			} else if code >= 300 && code < 400 {
				color = l.getColorName(l.theme.Status3xx)
			} else if code >= 400 && code < 500 {
				color = l.getColorName(l.theme.Status4xx)
			} else if code >= 500 {
				color = l.getColorName(l.theme.Status5xx)
			} else {
				color = l.getColorName(l.theme.Text)
			}
			statusText += fmt.Sprintf("[%d](fg:%s) %d responses\n", code, color, count)
		}
	}

	// 显示错误统计
	errorColor := l.getColorName(l.theme.Error)
	if errors > 0 {
		statusText += fmt.Sprintf("\n[Errors](fg:%s,mod:bold) : %d\n", errorColor, errors)
		errorRate := float64(errors) / float64(requests) * 100
		statusText += fmt.Sprintf("[Error Rate](fg:%s) : %.2f%%", errorColor, errorRate)
	}

	// 如果没有任何数据，显示提示信息
	warnColor := l.getColorName(l.theme.ReqSec)
	if len(statusCodes) == 0 && errors == 0 {
		statusText = fmt.Sprintf("[No responses yet...](fg:%s)", warnColor)
	} else if len(statusCodes) == 0 && errors > 0 {
		statusText += fmt.Sprintf("\n[All requests failed](fg:%s,mod:bold)", errorColor)
	}

	l.statusTable.Text = statusText

	// Update request chart
	if len(l.reqPerSecond) > 0 {
		data := make([]float64, len(l.reqPerSecond))
		labels := make([]string, len(l.reqPerSecond))
		for i, v := range l.reqPerSecond {
			data[i] = float64(v)
			labels[i] = fmt.Sprintf("%ds", i+1)
		}
		l.reqChart.Data = data
		l.reqChart.Labels = labels
	}

	// Update latency histogram
	if len(latencyPercentiles) > 0 {
		// 使用预计算的百分位
		buckets := []float64{50, 75, 90, 95, 99}
		data := make([]float64, len(buckets))
		labels := make([]string, len(buckets))

		for i, p := range buckets {
			if val, ok := latencyPercentiles[p]; ok {
				data[i] = val.Seconds() * 1000 // 转换为毫秒
			}
			labels[i] = fmt.Sprintf("p%.0f", p)
		}

		l.latencyChart.Data = data
		l.latencyChart.Labels = labels
	}
}

// UpdateEndpointStats updates statistics for a specific endpoint
func (l *LiveUI) UpdateEndpointStats(url string, requests int64, reqPerSec float64, avgLatency, minLatency, maxLatency time.Duration, errors int64) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.endpointStats[url] == nil {
		l.endpointStats[url] = &EndpointLiveStats{URL: url}
		l.multiEndpoint = true // 启用多端点模式
	}

	stats := l.endpointStats[url]
	stats.Requests = requests
	stats.ReqPerSec = reqPerSec
	stats.AvgLatency = avgLatency
	stats.MinLatency = minLatency
	stats.MaxLatency = maxLatency
	stats.Errors = errors
	stats.LastUpdate = time.Now()

	// 更新端点表格
	l.updateEndpointTable()
}

// updateEndpointTable updates the endpoint statistics table
func (l *LiveUI) updateEndpointTable() {
	if len(l.endpointStats) == 0 {
		return
	}

	// 表头
	rows := [][]string{
		{"Endpoint", "Req/s", "Avg", "Min", "Max", "Errors"},
	}

	// 按 URL 排序
	urls := make([]string, 0, len(l.endpointStats))
	for url := range l.endpointStats {
		urls = append(urls, url)
	}

	// 添加每个端点的数据
	for _, url := range urls {
		stats := l.endpointStats[url]
		errorRate := ""
		if stats.Requests > 0 {
			errorPct := float64(stats.Errors) / float64(stats.Requests) * 100
			if stats.Errors > 0 {
				errorRate = fmt.Sprintf("%d(%.1f%%)", stats.Errors, errorPct)
			} else {
				errorRate = "0"
			}
		}

		// 缩短 URL 显示
		displayURL := url
		if len(displayURL) > 40 {
			displayURL = displayURL[:37] + "..."
		}

		rows = append(rows, []string{
			displayURL,
			fmt.Sprintf("%.1f", stats.ReqPerSec),
			formatDurationShort(stats.AvgLatency),
			formatDurationShort(stats.MinLatency),
			formatDurationShort(stats.MaxLatency),
			errorRate,
		})
	}

	l.endpointTable.Rows = rows
}

// formatDurationShort formats duration in a short format
func formatDurationShort(d time.Duration) string {
	if d < time.Microsecond {
		return fmt.Sprintf("%.0fns", float64(d.Nanoseconds()))
	} else if d < time.Millisecond {
		return fmt.Sprintf("%.0fus", float64(d.Nanoseconds())/1000.0)
	} else if d < time.Second {
		return fmt.Sprintf("%.1fms", float64(d.Nanoseconds())/1000000.0)
	}
	return fmt.Sprintf("%.2fs", d.Seconds())
}

// Render renders the UI
func (l *LiveUI) Render() {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if l.multiEndpoint {
		// 多端点模式：显示端点表格
		ui.Render(l.progressGauge, l.statsTable, l.statusTable, l.reqChart, l.latencyChart, l.endpointTable, l.helpText)
	} else {
		// 单端点模式：不显示端点表格
		ui.Render(l.progressGauge, l.statsTable, l.statusTable, l.reqChart, l.latencyChart, l.helpText)
	}
}

// handleKeyEvents handles keyboard events
func (l *LiveUI) handleKeyEvents() {
	uiEvents := ui.PollEvents()
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				// 用户按下 'q' 或 Ctrl+C，触发停止
				close(l.stopChan)
				return
			}
		case <-l.stopChan:
			return
		}
	}
}

// StopChan returns the stop channel
func (l *LiveUI) StopChan() <-chan struct{} {
	return l.stopChan
}

// getColorName converts ui.Color to string name for text formatting
func (l *LiveUI) getColorName(color ui.Color) string {
	switch color {
	case ui.ColorBlack:
		return "black"
	case ui.ColorRed:
		return "red"
	case ui.ColorGreen:
		return "green"
	case ui.ColorYellow:
		return "yellow"
	case ui.ColorBlue:
		return "blue"
	case ui.ColorMagenta:
		return "magenta"
	case ui.ColorCyan:
		return "cyan"
	case ui.ColorWhite:
		return "white"
	default:
		// 支持 256 色模式，直接返回颜色代码
		// Paper主题使用的256色代码：
		// 25 (#005fd7), 64 (#008700), 160 (#d70000),
		// 166 (#df5f00), 172 (#d75f00), 238 (#444444), 244 (#808080)
		return fmt.Sprintf("%d", color)
	}
}

// Close closes the UI
func (l *LiveUI) Close() {
	// 确保 stopChan 被关闭
	select {
	case <-l.stopChan:
		// 已经关闭
	default:
		close(l.stopChan)
	}
	ui.Close()
}

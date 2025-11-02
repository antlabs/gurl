package benchmark

import (
	"fmt"
	"sync"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// LiveUI represents a live terminal UI for benchmark progress
type LiveUI struct {
	mu sync.RWMutex
	
	// Widgets
	progressGauge *widgets.Gauge
	statsTable    *widgets.Paragraph
	statusTable   *widgets.Paragraph
	reqChart      *widgets.BarChart
	latencyChart  *widgets.BarChart
	helpText      *widgets.Paragraph
	
	// Data
	totalRequests int64
	reqPerSecond  []int64
	statusCodes   map[int]int64
	latencies     []time.Duration
	
	startTime     time.Time
	duration      time.Duration
	
	// Control
	stopChan      chan struct{}
}

// NewLiveUI creates a new live UI
func NewLiveUI(duration time.Duration) (*LiveUI, error) {
	if err := ui.Init(); err != nil {
		return nil, err
	}
	
	liveUI := &LiveUI{
		statusCodes: make(map[int]int64),
		startTime:   time.Now(),
		duration:    duration,
		stopChan:    make(chan struct{}),
	}
	
	// Progress gauge
	liveUI.progressGauge = widgets.NewGauge()
	liveUI.progressGauge.Title = "Progress"
	liveUI.progressGauge.SetRect(0, 0, 100, 3)
	liveUI.progressGauge.BarColor = ui.ColorBlue
	liveUI.progressGauge.BorderStyle.Fg = ui.ColorWhite
	
	// Stats table
	liveUI.statsTable = widgets.NewParagraph()
	liveUI.statsTable.Title = "Stats for last sec"
	liveUI.statsTable.SetRect(0, 3, 50, 12)
	liveUI.statsTable.BorderStyle.Fg = ui.ColorWhite
	
	// Status code table
	liveUI.statusTable = widgets.NewParagraph()
	liveUI.statusTable.Title = "Status code distribution"
	liveUI.statusTable.SetRect(50, 3, 100, 12)
	liveUI.statusTable.BorderStyle.Fg = ui.ColorWhite
	
	// Request chart
	liveUI.reqChart = widgets.NewBarChart()
	liveUI.reqChart.Title = "Requests / past sec (auto)"
	liveUI.reqChart.SetRect(0, 12, 50, 30)
	liveUI.reqChart.BarWidth = 7  // 增加宽度以显示更多数字
	liveUI.reqChart.BarGap = 1    // 柱子间隔
	liveUI.reqChart.BarColors = []ui.Color{ui.ColorGreen}
	liveUI.reqChart.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorWhite)}
	liveUI.reqChart.NumStyles = []ui.Style{ui.NewStyle(ui.ColorYellow)}  // 改为黄色，更清晰
	liveUI.reqChart.BorderStyle.Fg = ui.ColorWhite
	
	// Latency chart
	liveUI.latencyChart = widgets.NewBarChart()
	liveUI.latencyChart.Title = "Response time histogram (ms)"
	liveUI.latencyChart.SetRect(50, 12, 100, 30)
	liveUI.latencyChart.BarWidth = 7  // 增加宽度
	liveUI.latencyChart.BarGap = 1    // 柱子间隔
	liveUI.latencyChart.BarColors = []ui.Color{ui.ColorYellow}
	liveUI.latencyChart.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorWhite)}
	liveUI.latencyChart.NumStyles = []ui.Style{ui.NewStyle(ui.ColorCyan)}  // 改为青色，更清晰
	liveUI.latencyChart.BorderStyle.Fg = ui.ColorWhite
	
	// Help text
	liveUI.helpText = widgets.NewParagraph()
	liveUI.helpText.Text = "[Press 'q' or 'Ctrl+C' to stop test early](fg:yellow,mod:bold)"
	liveUI.helpText.SetRect(0, 30, 100, 33)
	liveUI.helpText.Border = false
	liveUI.helpText.TextStyle.Fg = ui.ColorYellow
	
	// 启动键盘事件监听
	go liveUI.handleKeyEvents()
	
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
	
	// Update stats
	l.statsTable.Text = fmt.Sprintf(
		"[Requests](fg:white,mod:bold) : %d\n"+
		"[Slowest](fg:red) : %.4f secs\n"+
		"[Fastest](fg:green) : %.4f secs\n"+
		"[Average](fg:cyan) : %.4f secs\n"+
		"[Req/Sec](fg:yellow) : %.2f",
		requests,
		maxLatency.Seconds(),
		minLatency.Seconds(),
		avgLatency.Seconds(),
		float64(reqPerSec),
	)
	
	// Update status codes
	statusText := ""
	if len(statusCodes) > 0 {
		for code, count := range statusCodes {
			// 根据状态码类型设置颜色
			color := "green"
			if code >= 400 && code < 500 {
				color = "yellow" // 4xx 客户端错误
			} else if code >= 500 {
				color = "red" // 5xx 服务器错误
			}
			statusText += fmt.Sprintf("[%d](fg:%s) %d responses\n", code, color, count)
		}
	}
	
	// 显示错误统计
	if errors > 0 {
		statusText += fmt.Sprintf("\n[Errors](fg:red,mod:bold) : %d\n", errors)
		errorRate := float64(errors) / float64(requests) * 100
		statusText += fmt.Sprintf("[Error Rate](fg:red) : %.2f%%", errorRate)
	}
	
	// 如果没有任何数据，显示提示信息
	if len(statusCodes) == 0 && errors == 0 {
		statusText = "[No responses yet...](fg:yellow)"
	} else if len(statusCodes) == 0 && errors > 0 {
		statusText += "\n[All requests failed](fg:red,mod:bold)"
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

// Render renders the UI
func (l *LiveUI) Render() {
	l.mu.RLock()
	defer l.mu.RUnlock()
	
	ui.Render(l.progressGauge, l.statsTable, l.statusTable, l.reqChart, l.latencyChart, l.helpText)
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

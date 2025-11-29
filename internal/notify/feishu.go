package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/antlabs/gurl/internal/batch"
	"github.com/antlabs/gurl/internal/config"
)

// Notifier defines interface for sending batch notifications.
type Notifier interface {
	NotifyBatch(result *batch.BatchResult, report string) error
}

// NewFromConfig creates a Notifier from batch notifier config.
func NewFromConfig(cfg *config.NotifierConfig) Notifier {
	if cfg == nil {
		return nil
	}

	// default Enabled to true when not explicitly set but webhook exists
	if !cfg.Enabled && cfg.FeishuWebhook != "" {
		cfg.Enabled = true
	}

	switch cfg.Type {
	case "feishu", "feishu_webhook", "feishu-webhook", "":
		// default type is feishu when webhook is provided
		if cfg.FeishuWebhook == "" {
			return nil
		}
		return &FeishuNotifier{cfg: cfg}
	default:
		return nil
	}
}

// FeishuNotifier sends batch result notifications to Feishu via incoming webhook.
type FeishuNotifier struct {
	cfg *config.NotifierConfig
}

// feishuCardMessage is a minimal card message payload for Feishu.
type feishuCardMessage struct {
	MsgType string         `json:"msg_type"`
	Card    feishuCardBody `json:"card"`
}

type feishuCardBody struct {
	Config   feishuCardConfig `json:"config"`
	Header   feishuCardHeader `json:"header"`
	Elements []any            `json:"elements"`
}

type feishuCardConfig struct {
	WideScreenMode bool `json:"wide_screen_mode"`
}

type feishuCardHeader struct {
	Title feishuCardTitle `json:"title"`
}

type feishuCardTitle struct {
	Tag     string `json:"tag"`
	Content string `json:"content"`
}

// NotifyBatch implements Notifier.
func (n *FeishuNotifier) NotifyBatch(result *batch.BatchResult, report string) error {
	if n == nil || n.cfg == nil {
		return nil
	}
	if !n.cfg.Enabled {
		return nil
	}

	// Optionally only notify when there are failures.
	if n.cfg.OnlyOnFail {
		hasFailure := false
		for _, t := range result.Tests {
			if t.Error != nil {
				hasFailure = true
				break
			}
			if t.Stats != nil && len(t.Stats.GetErrors()) > 0 {
				hasFailure = true
				break
			}
		}
		if !hasFailure {
			return nil
		}
	}

	title := n.cfg.Title
	if title == "" {
		// Basic default title with success rate
		title = fmt.Sprintf("gurl batch %.1f%% passed", result.SuccessRate)
	}

	// For now send the plain text report inside a markdown element.
	text := fmt.Sprintf("```\n%s\n```", report)

	payload := feishuCardMessage{
		MsgType: "interactive",
		Card: feishuCardBody{
			Config: feishuCardConfig{WideScreenMode: true},
			Header: feishuCardHeader{
				Title: feishuCardTitle{Tag: "plain_text", Content: title},
			},
			Elements: []any{
				map[string]any{
					"tag":     "markdown",
					"content": text,
				},
			},
		},
	}

	buf, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest(http.MethodPost, n.cfg.FeishuWebhook, bytes.NewReader(buf))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("feishu webhook returned status %d", resp.StatusCode)
	}

	return nil
}

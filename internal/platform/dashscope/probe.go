package dashscope

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/NERVEbing/model-scout/internal/platform"
)

type probeRequest struct {
	Model    string `json:"model"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
	MaxTokens int `json:"max_tokens"`
}

func (p *Platform) Probe(ctx context.Context, model platform.Model) platform.ProbeResult {
	request := probeRequest{
		Model: model.ID,
		Messages: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{{Role: "user", Content: "ping"}},
		MaxTokens: 1,
	}
	payload, err := json.Marshal(request)
	if err != nil {
		return platform.ProbeResult{
			Platform:  p.Name(),
			Model:     model.ID,
			Status:    "error",
			Available: false,
			Reason:    err.Error(),
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.client.BaseURL+"/chat/completions", bytes.NewReader(payload))
	if err != nil {
		return platform.ProbeResult{
			Platform:  p.Name(),
			Model:     model.ID,
			Status:    "error",
			Available: false,
			Reason:    err.Error(),
		}
	}
	req.Header.Set("Authorization", "Bearer "+p.client.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.HTTPClient.Do(req)
	if err != nil {
		return platform.ProbeResult{
			Platform:  p.Name(),
			Model:     model.ID,
			Status:    "error",
			Available: false,
			Reason:    err.Error(),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return platform.ProbeResult{
			Platform:     p.Name(),
			Model:        model.ID,
			Status:       "ok",
			Available:    true,
			Capabilities: []string{"chat"},
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return platform.ProbeResult{
			Platform:  p.Name(),
			Model:     model.ID,
			Status:    "error",
			Available: false,
			Reason:    err.Error(),
		}
	}

	reason := strings.TrimSpace(string(body))
	if reason != "" {
		reason = fmt.Sprintf("%s: %s", resp.Status, reason)
	} else {
		reason = resp.Status
	}

	return platform.ProbeResult{
		Platform:  p.Name(),
		Model:     model.ID,
		Status:    "fail",
		Available: false,
		Reason:    reason,
	}
}

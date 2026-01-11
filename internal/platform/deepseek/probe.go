package deepseek

import (
	"bytes"
	"context"
	"encoding/json"
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

type errorEnvelope struct {
	Error struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
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

	var envelope errorEnvelope
	if err := json.Unmarshal(body, &envelope); err != nil {
		return platform.ProbeResult{
			Platform:  p.Name(),
			Model:     model.ID,
			Status:    "fail",
			Available: false,
			Reason:    resp.Status,
		}
	}

	errorType := envelope.Error.Type
	reason := envelope.Error.Message
	if reason == "" {
		reason = resp.Status
	}

	status := "fail"
	if strings.Contains(errorType, "AccessDenied") || strings.Contains(errorType, "Model.AccessDenied") {
		status = "denied"
	} else if strings.Contains(errorType, "NotSupported") {
		status = "unsupported"
	}

	return platform.ProbeResult{
		Platform:  p.Name(),
		Model:     model.ID,
		Status:    status,
		Available: false,
		Reason:    reason,
	}
}

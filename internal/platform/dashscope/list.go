package dashscope

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/NERVEbing/model-scout/internal/platform"
)

type listResponse struct {
	Data []struct {
		ID string `json:"id"`
	} `json:"data"`
}

func (p *Platform) ListModels(ctx context.Context) ([]platform.Model, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.client.BaseURL+"/models", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+p.client.APIKey)

	resp, err := p.client.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("dashscope list models failed: %s", resp.Status)
	}

	var payload listResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	models := make([]platform.Model, 0, len(payload.Data))
	for _, item := range payload.Data {
		if item.ID == "" {
			continue
		}
		models = append(models, platform.Model{ID: item.ID})
	}
	return models, nil
}

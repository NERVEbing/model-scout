package dashscope

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListModels(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/models" {
			http.NotFound(w, r)
			return
		}
		if got := r.Header.Get("Authorization"); got != "Bearer token" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":"qwen-plus"},{"id":"qwen-turbo"}]}`))
	}))
	t.Cleanup(server.Close)

	platform := &Platform{
		client: &Client{
			BaseURL:    server.URL,
			APIKey:     "token",
			HTTPClient: server.Client(),
		},
	}

	models, err := platform.ListModels(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(models) != 2 {
		t.Fatalf("expected 2 models, got %d", len(models))
	}
	if models[0].ID != "qwen-plus" || models[1].ID != "qwen-turbo" {
		t.Fatalf("unexpected models: %#v", models)
	}
}

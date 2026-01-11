package dashscope

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
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

type errorBody struct{}

func (errorBody) Read([]byte) (int, error) {
	return 0, errors.New("read failed")
}

func (errorBody) Close() error {
	return nil
}

type errorTransport struct{}

func (errorTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusUnauthorized,
		Status:     "401 Unauthorized",
		Body:       errorBody{},
		Header:     make(http.Header),
	}, nil
}

func TestListModelsReadError(t *testing.T) {
	client := &Client{
		BaseURL: "https://example.com",
		APIKey:  "token",
		HTTPClient: &http.Client{
			Transport: errorTransport{},
		},
	}
	platform := &Platform{client: client}

	_, err := platform.ListModels(context.Background())
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "read body") {
		t.Fatalf("expected read body error, got %v", err)
	}
}

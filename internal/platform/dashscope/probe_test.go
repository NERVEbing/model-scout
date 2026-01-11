package dashscope

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/NERVEbing/model-scout/internal/platform"
)

func TestProbeOK(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chat/completions" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	platformImpl := &Platform{
		client: &Client{
			BaseURL:    server.URL,
			APIKey:     "token",
			HTTPClient: server.Client(),
		},
	}

	result := platformImpl.Probe(context.Background(), platform.Model{ID: "qwen-plus"})
	if result.Status != "ok" || !result.Available {
		t.Fatalf("expected ok/available, got status=%s available=%t", result.Status, result.Available)
	}
}

func TestProbeFail(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chat/completions" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte("no access"))
	}))
	t.Cleanup(server.Close)

	platformImpl := &Platform{
		client: &Client{
			BaseURL:    server.URL,
			APIKey:     "token",
			HTTPClient: server.Client(),
		},
	}

	result := platformImpl.Probe(context.Background(), platform.Model{ID: "qwen-plus"})
	if result.Status != "fail" || result.Available {
		t.Fatalf("expected fail/unavailable, got status=%s available=%t", result.Status, result.Available)
	}
	if !strings.Contains(result.Reason, "403") || !strings.Contains(result.Reason, "no access") {
		t.Fatalf("unexpected reason: %s", result.Reason)
	}
}

package cli

import (
	"testing"

	"github.com/NERVEbing/model-scout/internal/platform"
)

func TestParseFilters(t *testing.T) {
	t.Run("valid eq", func(t *testing.T) {
		filters, err := parseFilters([]string{"status=ok,active"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(filters) != 1 {
			t.Fatalf("expected 1 filter, got %d", len(filters))
		}
		filter := filters[0]
		if filter.key != "status" {
			t.Fatalf("expected key status, got %q", filter.key)
		}
		if filter.op != filterEq {
			t.Fatalf("expected eq op")
		}
		if len(filter.values) != 2 || filter.values[0] != "ok" || filter.values[1] != "active" {
			t.Fatalf("unexpected values: %#v", filter.values)
		}
	})

	t.Run("valid not eq", func(t *testing.T) {
		filters, err := parseFilters([]string{"model!=qwen-plus"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(filters) != 1 {
			t.Fatalf("expected 1 filter, got %d", len(filters))
		}
		if filters[0].op != filterNotEq {
			t.Fatalf("expected not eq op")
		}
	})

	t.Run("invalid format", func(t *testing.T) {
		_, err := parseFilters([]string{"status"})
		if err == nil {
			t.Fatalf("expected error for invalid format")
		}
	})

	t.Run("unsupported key", func(t *testing.T) {
		_, err := parseFilters([]string{"unknown=ok"})
		if err == nil {
			t.Fatalf("expected error for unsupported key")
		}
	})

	t.Run("invalid boolean", func(t *testing.T) {
		_, err := parseFilters([]string{"available=maybe"})
		if err == nil {
			t.Fatalf("expected error for invalid boolean")
		}
	})
}

func TestApplyFilters(t *testing.T) {
	results := []platform.ProbeResult{
		{Platform: "dashscope", Model: "qwen-plus", Status: "ok", Available: true},
		{Platform: "dashscope", Model: "qwen-turbo", Status: "active", Available: true},
		{Platform: "dashscope", Model: "qwen-mini", Status: "fail", Available: false},
		{Platform: "other", Model: "other-model", Status: "ok", Available: true},
	}

	t.Run("filter status", func(t *testing.T) {
		filters, err := parseFilters([]string{"status=ok"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		filtered := applyFilters(results, filters)
		if len(filtered) != 2 {
			t.Fatalf("expected 2 results, got %d", len(filtered))
		}
	})

	t.Run("filter available", func(t *testing.T) {
		filters, err := parseFilters([]string{"available=true"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		filtered := applyFilters(results, filters)
		if len(filtered) != 3 {
			t.Fatalf("expected 3 results, got %d", len(filtered))
		}
	})

	t.Run("filter not eq", func(t *testing.T) {
		filters, err := parseFilters([]string{"platform!=dashscope"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		filtered := applyFilters(results, filters)
		if len(filtered) != 1 {
			t.Fatalf("expected 1 result, got %d", len(filtered))
		}
		if filtered[0].Platform != "other" {
			t.Fatalf("unexpected platform: %s", filtered[0].Platform)
		}
	})

	t.Run("filter or within key", func(t *testing.T) {
		filters, err := parseFilters([]string{"status=ok,active"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		filtered := applyFilters(results, filters)
		if len(filtered) != 3 {
			t.Fatalf("expected 3 results, got %d", len(filtered))
		}
	})

	t.Run("filter multiple keys", func(t *testing.T) {
		filters, err := parseFilters([]string{"platform=dashscope", "status=ok"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		filtered := applyFilters(results, filters)
		if len(filtered) != 1 {
			t.Fatalf("expected 1 result, got %d", len(filtered))
		}
		if filtered[0].Model != "qwen-plus" {
			t.Fatalf("unexpected model: %s", filtered[0].Model)
		}
	})
}

func TestDefaultKeyEnv(t *testing.T) {
	t.Run("dashscope default", func(t *testing.T) {
		env, err := defaultKeyEnv("dashscope")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if env != defaultDashscopeKeyEnv {
			t.Fatalf("expected %s, got %q", defaultDashscopeKeyEnv, env)
		}
	})

	t.Run("deepseek default", func(t *testing.T) {
		env, err := defaultKeyEnv("deepseek")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if env != "DEEPSEEK_API_KEY" {
			t.Fatalf("expected DEEPSEEK_API_KEY, got %q", env)
		}
	})

	t.Run("unsupported platform", func(t *testing.T) {
		_, err := defaultKeyEnv("unknown")
		if err == nil {
			t.Fatalf("expected error for unsupported platform")
		}
	})
}

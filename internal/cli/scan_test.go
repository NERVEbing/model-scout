package cli

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/NERVEbing/model-scout/internal/platform"
)

type fakePlatform struct{}

func (p *fakePlatform) Name() string {
	return "fake"
}

func (p *fakePlatform) ListModels(_ context.Context) ([]platform.Model, error) {
	return []platform.Model{
		{ID: "ok-model"},
		{ID: "skip-model"},
		{ID: "fail-model"},
	}, nil
}

func (p *fakePlatform) Probe(_ context.Context, model platform.Model) platform.ProbeResult {
	switch model.ID {
	case "ok-model", "skip-model":
		return platform.ProbeResult{
			Platform:  p.Name(),
			Model:     model.ID,
			Status:    "ok",
			Available: true,
		}
	default:
		return platform.ProbeResult{
			Platform:  p.Name(),
			Model:     model.ID,
			Status:    "fail",
			Available: false,
			Reason:    "failed",
		}
	}
}

func TestRunOutputsFilteredResults(t *testing.T) {
	t.Helper()

	prevFactory := platformFactory
	platformFactory = func(_ string, _ string, _ time.Duration) (platform.Platform, error) {
		return &fakePlatform{}, nil
	}
	t.Cleanup(func() {
		platformFactory = prevFactory
	})

	outputPath := filepath.Join(t.TempDir(), "out.json")
	t.Setenv("DEEPSEEK_API_KEY", "token")
	args := []string{
		"--platform", "deepseek",
		"--output-file", outputPath,
		"--out", "json",
		"--exclude", "skip",
		"--filter", "status=ok",
		"--only-ok",
	}
	if err := Run(args); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}

	var results []platform.ProbeResult
	if err := json.Unmarshal(data, &results); err != nil {
		t.Fatalf("unmarshal output: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Model != "ok-model" {
		t.Fatalf("unexpected model: %s", results[0].Model)
	}
	if results[0].Platform != "fake" {
		t.Fatalf("unexpected platform: %s", results[0].Platform)
	}
}

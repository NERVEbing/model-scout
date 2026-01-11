package scout

import (
	"context"
	"sync"
	"testing"

	"github.com/NERVEbing/model-scout/internal/platform"
)

type fakePlatform struct {
	mu       sync.Mutex
	probed   map[string]bool
	toReturn []platform.Model
}

func (f *fakePlatform) Name() string {
	return "fake"
}

func (f *fakePlatform) ListModels(ctx context.Context) ([]platform.Model, error) {
	return f.toReturn, nil
}

func (f *fakePlatform) Probe(ctx context.Context, model platform.Model) platform.ProbeResult {
	f.mu.Lock()
	if f.probed == nil {
		f.probed = make(map[string]bool)
	}
	f.probed[model.ID] = true
	f.mu.Unlock()

	return platform.ProbeResult{
		Platform:  f.Name(),
		Model:     model.ID,
		Status:    "ok",
		Available: true,
	}
}

func TestEngineScanFilters(t *testing.T) {
	fake := &fakePlatform{
		toReturn: []platform.Model{
			{ID: "model-image"},
			{ID: "text-1"},
			{ID: "tts-foo"},
			{ID: "chat-model"},
		},
	}

	engine := Engine{Platform: fake, Workers: 2}
	results, err := engine.Scan(context.Background(), []string{"chat"})
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Model != "text-1" {
		t.Fatalf("expected text-1, got %s", results[0].Model)
	}
}

package scout

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

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

type cancelPlatform struct {
	started  chan struct{}
	toReturn []platform.Model
}

func (c *cancelPlatform) Name() string {
	return "cancel"
}

func (c *cancelPlatform) ListModels(ctx context.Context) ([]platform.Model, error) {
	return c.toReturn, nil
}

func (c *cancelPlatform) Probe(ctx context.Context, model platform.Model) platform.ProbeResult {
	select {
	case c.started <- struct{}{}:
	default:
	}
	<-ctx.Done()
	return platform.ProbeResult{
		Platform:  c.Name(),
		Model:     model.ID,
		Status:    "error",
		Available: false,
		Reason:    ctx.Err().Error(),
	}
}

func TestEngineScanRespectsContextCancel(t *testing.T) {
	cancelSignal := make(chan struct{}, 1)
	fake := &cancelPlatform{
		started: cancelSignal,
		toReturn: []platform.Model{
			{ID: "model-a"},
			{ID: "model-b"},
		},
	}
	engine := Engine{Platform: fake, Workers: 2}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct{})
	var results []platform.ProbeResult
	var err error
	go func() {
		results, err = engine.Scan(ctx, nil)
		close(done)
	}()

	select {
	case <-cancelSignal:
	case <-time.After(1 * time.Second):
		t.Fatalf("expected probe to start before cancel")
	}

	cancel()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatalf("expected scan to return after cancel")
	}

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context canceled error, got %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected no results, got %d", len(results))
	}
}

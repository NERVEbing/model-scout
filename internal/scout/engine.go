package scout

import (
	"context"
	"fmt"
	"sync"

	"github.com/NERVEbing/model-scout/internal/platform"
)

type Engine struct {
	Platform platform.Platform
	Workers  int
}

func (e Engine) Scan(ctx context.Context, excludes []string) ([]platform.ProbeResult, error) {
	if e.Platform == nil {
		return nil, fmt.Errorf("platform is required")
	}
	workers := e.Workers
	if workers <= 0 {
		workers = 1
	}

	models, err := e.Platform.ListModels(ctx)
	if err != nil {
		return nil, err
	}
	filtered := make([]platform.Model, 0, len(models))
	for _, model := range models {
		if ShouldSkip(model.ID, excludes) {
			continue
		}
		filtered = append(filtered, model)
	}

	jobs := make(chan platform.Model)
	results := make(chan platform.ProbeResult)

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for model := range jobs {
				results <- e.Platform.Probe(ctx, model)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	go func() {
		for _, model := range filtered {
			jobs <- model
		}
		close(jobs)
	}()

	collected := make([]platform.ProbeResult, 0, len(filtered))
	for result := range results {
		collected = append(collected, result)
	}

	return collected, nil
}

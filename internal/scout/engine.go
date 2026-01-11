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
	ctxDone := ctx.Done()

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Go(func() {
			for {
				select {
				case <-ctxDone:
					return
				case model, ok := <-jobs:
					if !ok {
						return
					}
					result := e.Platform.Probe(ctx, model)
					select {
					case <-ctxDone:
						return
					case results <- result:
					}
				}
			}
		})
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	go func() {
		defer close(jobs)
		for _, model := range filtered {
			select {
			case <-ctxDone:
				return
			case jobs <- model:
			}
		}
	}()

	collected := make([]platform.ProbeResult, 0, len(filtered))
	canceled := false
	for {
		select {
		case result, ok := <-results:
			if !ok {
				if canceled && ctx.Err() != nil {
					return collected, ctx.Err()
				}
				return collected, nil
			}
			collected = append(collected, result)
		case <-ctxDone:
			canceled = true
		}
	}
}

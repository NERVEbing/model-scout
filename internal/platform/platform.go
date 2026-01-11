package platform

import "context"

type Platform interface {
	Name() string
	ListModels(ctx context.Context) ([]Model, error)
	Probe(ctx context.Context, model Model) ProbeResult
}

type Model struct {
	ID   string
	Meta map[string]string
}

type ProbeResult struct {
	Platform     string            `json:"platform" yaml:"platform"`
	Model        string            `json:"model" yaml:"model"`
	Status       string            `json:"status" yaml:"status"`
	Available    bool              `json:"available" yaml:"available"`
	Reason       string            `json:"reason,omitempty" yaml:"reason,omitempty"`
	Capabilities []string          `json:"capabilities,omitempty" yaml:"capabilities,omitempty"`
	Meta         map[string]string `json:"meta,omitempty" yaml:"meta,omitempty"`
}

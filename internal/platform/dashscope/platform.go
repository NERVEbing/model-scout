package dashscope

import "time"

type Platform struct {
	client *Client
}

func NewPlatform(apiKey string, timeout time.Duration) *Platform {
	return &Platform{client: NewClient(apiKey, timeout)}
}

func (p *Platform) Name() string {
	return "dashscope"
}

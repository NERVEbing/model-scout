package dashscope

import (
	"net/http"
	"time"
)

type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

func NewClient(apiKey string, timeout time.Duration) *Client {
	return &Client{
		BaseURL: "https://dashscope.aliyuncs.com/compatible-mode/v1",
		APIKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: timeout,
		},
	}
}

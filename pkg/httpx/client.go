package httpx

import (
	"net/http"
	"time"
)

type ClientOption func(*http.Client)

func NewClient(opts ...ClientOption) *http.Client {
	c := &http.Client{}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func ClientOptionWithTimeout(d time.Duration) ClientOption {
	return func(c *http.Client) {
		c.Timeout = d
	}
}

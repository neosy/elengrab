package httpx

import (
	"net/http"
)

// NewClient creates a new http.Client and applies any number of ClientOptions.
// opts: variadic list of ClientOption functions to customize the client.
func NewClient(opts ...ClientOption) *http.Client {
	c := &http.Client{}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

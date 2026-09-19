package fasthttpx

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

// GetHTML performs an HTTP GET request using fasthttp and returns
// the response body as a safe-to-use byte slice.
func GetHTML(url string, opts ...ClientOption) ([]byte, error) {
	var req fasthttp.Request
	var resp fasthttp.Response

	// Set target URL
	req.SetRequestURI(url)
	req.Header.SetMethod(fasthttp.MethodGet)

	// Set a realistic browser User-Agent to match curl / browser behavior
	req.Header.SetUserAgent(
		LinuxUserAgent,
	)

	// Basic headers similar to curl defaults
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	// Create a client with increased read buffer size
	client := NewClient(opts...)

	// Execute the HTTP request and follow redirects
	execRequest := func(client *fasthttp.Client) error {
		if client.ReadTimeout > 0 {
			req.SetTimeout(client.ReadTimeout)
		}

		return client.DoRedirects(&req, &resp, 10)
	}
	if err := execRequest(client); err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	// Ensure successful HTTP response
	if resp.StatusCode() != fasthttp.StatusOK {
		return nil, fmt.Errorf(
			"unexpected status code: %d, server: %s, retry-after: %s, content-type: %s, url: %s",
			resp.StatusCode(),
			resp.Header.Peek("Server"),
			resp.Header.Peek("Retry-After"),
			resp.Header.Peek("Content-Type"),
			req.URI().String(),
		)
	}

	// Copy body to a new slice to prevent buffer reuse issues
	body := make([]byte, len(resp.Body()))
	copy(body, resp.Body())

	return body, nil
}

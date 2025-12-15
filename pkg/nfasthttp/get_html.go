package nfasthttp

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

// GetHTML performs an HTTP GET request using fasthttp and returns
// the response body as a safe-to-use byte slice.
// bufferSize:  64 * 1024 (for youtube)
func GetHTML(url string, bufferSize int) ([]byte, error) {
	var req fasthttp.Request
	var resp fasthttp.Response

	// Set target URL
	req.SetRequestURI(url)
	req.Header.SetMethod(fasthttp.MethodGet)

	// Set a realistic browser User-Agent to match curl / browser behavior
	req.Header.SetUserAgent(
		linuxUserAgent,
	)

	// Basic headers similar to curl defaults
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	client := &fasthttp.Client{
		ReadBufferSize: bufferSize,
	}

	// Execute the HTTP request
	if err := client.Do(&req, &resp); err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	// Ensure successful HTTP response
	if resp.StatusCode() != fasthttp.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	// Copy body to a new slice to prevent buffer reuse issues
	body := make([]byte, len(resp.Body()))
	copy(body, resp.Body())

	return body, nil
}

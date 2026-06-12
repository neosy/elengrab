package fasthttpx

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

// GetImageFormatByHeader returns image format based on Content-Type header
func GetImageFormatByHeader(resp *fasthttp.Response) (string, error) {
	contentType := string(resp.Header.ContentType())
	if contentType == "" {
		return "", fmt.Errorf("Content-Type header not found")
	}

	switch contentType {
	case "image/jpeg":
		return "jpg", nil
	case "image/png":
		return "png", nil
	case "image/webp":
		return "webp", nil
	default:
		return "", fmt.Errorf("unsupported content type: %s", contentType)
	}
}

// GetImage downloads an image from the given URL using fasthttp
// and returns a safe-to-use copy of the image bytes along with its format.
func GetImage(url string, opts ...ClientOption) ([]byte, string, error) {
	var req fasthttp.Request
	var resp fasthttp.Response

	// Set request URL and method
	req.SetRequestURI(url)
	req.Header.SetMethod(fasthttp.MethodGet)

	// Set a realistic User-Agent to emulate a browser
	req.Header.SetUserAgent(LinuxUserAgent)

	// Create a client with increased read buffer size
	client := NewClient(opts...)

	// Execute the HTTP request
	execRequest := func(client *fasthttp.Client) error {
		if client.ReadTimeout > 0 {
			return client.DoTimeout(&req, &resp, client.ReadTimeout)
		}
		return client.Do(&req, &resp)
	}
	if err := execRequest(client); err != nil {
		return nil, "", fmt.Errorf("get image error: %v", err)
	}

	// Check HTTP status
	if resp.StatusCode() != fasthttp.StatusOK {
		return nil, "", fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	// Determine image format from Content-Type header
	format, err := GetImageFormatByHeader(&resp)
	if err != nil {
		return nil, "", fmt.Errorf("get image format error: %v", err)
	}

	// Copy body to a new slice to prevent buffer reuse issues
	body := make([]byte, len(resp.Body()))
	copy(body, resp.Body())

	return body, format, nil
}

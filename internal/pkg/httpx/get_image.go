package httpx

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// GetImage downloads an image from the given URL using net/http
// and returns a safe-to-use copy of the image bytes along with its format.
//
// The function accepts optional parameters (opts) to configure the request.
// Currently supported option types are:
//   - GetOptions   : general options for image retrieval (e.g., limit, etc.)
//     Only one GetOptions instance can be passed; later ones overwrite earlier ones.
//   - ClientOption : client-specific options (e.g., ClientOptionWithTimeout())
//     Multiple ClientOption instances can be passed and are all applied.//
func GetImage(ctx context.Context, url string, opts ...any) ([]byte, string, error) {
	var (
		getOpts    GetOptions
		clientOpts []ClientOption
	)

	for _, opt := range opts {
		switch v := opt.(type) {
		case GetOptions:
			getOpts = v
		case ClientOption:
			clientOpts = append(clientOpts, v)
		default:
			panic("unsupported option type")
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, "", err
	}

	// Set a realistic User-Agent
	req.Header.Set("User-Agent", LinuxUserAgent)

	client := NewClient(clientOpts...)

	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	maxSize := LimitImageDefault
	if getOpts.Limit != 0 {
		maxSize = getOpts.Limit
	}

	data, err := io.ReadAll(io.LimitReader(resp.Body, maxSize+1))
	if err != nil {
		return nil, "", err
	}

	if int64(len(data)) > maxSize {
		return nil, "", errors.New("image exceeds maximum allowed size")
	}

	// Determine image format from Content-Type header
	ext, err := GetImageExtFromHeader(resp)
	if err != nil {
		return nil, "", fmt.Errorf("get image format error: %v", err)
	}

	return data, ext, nil
}

// GetImageExtFromHeader returns image format based on Content-Type header
func GetImageExtFromHeader(resp *http.Response) (string, error) {
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		return "", fmt.Errorf("Content-Type header not found")
	}

	mediaType := ExtensionByContentType(contentType)
	if mediaType == "" {
		return "", fmt.Errorf("unsupported content type: %s", contentType)
	}

	return mediaType, nil
}

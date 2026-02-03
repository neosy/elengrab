package httpx

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/neosy/elengrab/pkg/nfasthttp"
)

// GetHTML performs an HTTP GET request using net/http and returns
// the response body as a safe-to-use byte slice.
//
// The function accepts optional parameters (opts) to configure the request.
// Currently supported option types are:
//   - GetOptions   : general options for retrieval (e.g., limit, etc.)
//     Only one GetOptions instance can be passed; later ones overwrite earlier ones.
//   - ClientOption : client-specific options (e.g., ClientOptionWithTimeout())
//     Multiple ClientOption instances can be passed and are all applied.//
func GetHTML(ctx context.Context, url string, opts ...any) ([]byte, error) {
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
		return nil, err
	}

	// Set a realistic User-Agent
	req.Header.Set("User-Agent", nfasthttp.LinuxUserAgent)

	// Basic headers similar to curl defaults
	req.Header.Set("Accept", "*/*")
	acceptLanguage := acceptLanguageDefault
	if getOpts.AcceptLanguage != "" {
		acceptLanguage = getOpts.AcceptLanguage
	}
	req.Header.Set("Accept-Language", acceptLanguage)

	client := NewClient(clientOpts...)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	maxSize := LimitHTMLDefault
	if getOpts.Limit != 0 {
		maxSize = getOpts.Limit
	}

	data, err := io.ReadAll(io.LimitReader(resp.Body, maxSize+1))
	if err != nil {
		return nil, err
	}

	if int64(len(data)) > maxSize {
		return nil, errors.New("html exceeds maximum allowed size")
	}

	return data, nil
}

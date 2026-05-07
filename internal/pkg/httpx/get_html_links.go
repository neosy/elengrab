package httpx

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/html"
)

// GetLinksInHead downloads the HTML from the given URL and extracts all <link> tags
// that are inside the <head> section. The parsing stops as soon as </head> or <body>
// is encountered, so the rest of the HTML is not read.
//
// The function accepts optional parameters (opts) to configure the request.
// Currently supported option types are:
//   - MethodGetOptions   : general options for retrieval (e.g., Limit, AcceptLanguage).
//     Only one MethodGetOptions instance can be passed; later ones overwrite earlier ones.
//   - ClientOption : client-specific options (e.g., ClientOptionWithTimeout()).
//     Multiple ClientOption instances can be passed and are all applied.
func GetLinksInHead(ctx context.Context, url string, opts ...any) ([][]html.Attribute, error) {
	var (
		getOpts    MethodGetOptions
		clientOpts []ClientOption
	)

	for _, opt := range opts {
		switch v := opt.(type) {
		case MethodGetOptions:
			getOpts = v
		case ClientOption:
			clientOpts = append(clientOpts, v)
		default:
			return nil, fmt.Errorf("unsupported option type %T", opt)
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// Set a realistic User-Agent
	req.Header.Set("User-Agent", LinuxUserAgent)

	// Basic headers similar to curl defaults
	req.Header.Set("Accept", "text/html")
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
	if (!getOpts.IgnoreStatusCode || resp == nil) && resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var (
		links     [][]html.Attribute
		tokenizer = html.NewTokenizer(resp.Body)
		inHead    = false
	)

	for {
		tt := tokenizer.Next()
		switch tt {
		case html.ErrorToken:
			if tokenizer.Err() == io.EOF {
				return links, nil
			}
			return nil, tokenizer.Err()
		case html.StartTagToken, html.SelfClosingTagToken:
			t := tokenizer.Token()
			switch t.Data {
			case "head":
				inHead = true
			case "body":
				// Beginning <body> — we are no longer looking for <link>
				return links, nil
			case "link":
				if inHead {
					links = append(links, t.Attr)
				}
			}
		case html.EndTagToken:
			t := tokenizer.Token()
			// If the end is </head> - we are no longer looking for <link>
			if t.Data == "head" {
				return links, nil
			}
		}
	}
}

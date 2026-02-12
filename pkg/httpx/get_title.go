package httpx

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
)

func GetTitle(ctx context.Context, url string, opts ...any) (string, error) {
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
			return "", fmt.Errorf("unsupported option type %T", opt)
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
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
		return "", err
	}
	defer resp.Body.Close()

	// Limit size to avoid excessive memory usage.
	maxSize := LimitHTMLTitleDefault
	if getOpts.Limit != 0 {
		maxSize = getOpts.Limit
	}

	// Create a reader for the response body with the appropriate character set
	// based on the Content-Type header.
	reader, _ := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))

	// Limit the body size to avoid excessive memory usage.
	body := io.LimitReader(reader, maxSize)

	// Check HTTP status
	if (!getOpts.IgnoreStatusCode || resp == nil) && resp.StatusCode >= 300 {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var (
		title     string
		tokenizer = html.NewTokenizer(body)
		inHead    = false
	)

	for {
		tt := tokenizer.Next()
		switch tt {
		case html.ErrorToken:
			if tokenizer.Err() == io.EOF {
				return title, nil
			}
			return "", tokenizer.Err()
		case html.StartTagToken, html.SelfClosingTagToken:
			t := tokenizer.Token()
			switch t.Data {
			case "head":
				inHead = true
			case "body":
				// Beginning <body> — we are no longer looking for <title>
				return "", nil
			case "title":
				if inHead {
					if tokenizer.Next() == html.TextToken {
						title = tokenizer.Token().Data
						return strings.TrimSpace(title), nil
					}
				}
			}
		case html.EndTagToken:
			t := tokenizer.Token()
			// If the end is </head> - we are no longer looking for <title>
			if t.Data == "head" {
				return "", nil
			}
		}
	}
}

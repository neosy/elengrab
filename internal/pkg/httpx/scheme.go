package httpx

import "net/url"

// SchemeFromURL parses the given URL and returns its scheme (e.g. "http", "https").
// If the URL is invalid or has no scheme, it returns an empty string.
func SchemeFromURL(rawURL string) string {
	if rawURL == "" {
		return ""
	}

	u, err := url.Parse(rawURL)
	if err == nil && u.Scheme != "" {
		return u.Scheme
	}

	return ""
}

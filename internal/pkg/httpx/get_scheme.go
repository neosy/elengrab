package httpx

import "net/url"

// GetSchemeFromURL parses a URL and returns the scheme part of the URL.
func GetSchemeFromURL(rawURL string) string {
	if rawURL == "" {
		return ""
	}

	u, err := url.Parse(rawURL)
	if err == nil && u.Scheme != "" {
		return u.Scheme
	}

	return ""
}

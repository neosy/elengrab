package httpx

import (
	"fmt"
	"net/url"
)

// GetBaseURL returns the base URL (scheme + host) from a full URL string.
// If parsing fails, it returns an empty string.
func GetBaseURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return ""
	}
	return fmt.Sprintf("%s://%s", u.Scheme, u.Host)
}

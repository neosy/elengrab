package linkweb

import (
	"net/url"
	"strings"
)

func (u *linkWeb) isAbsBaseShortURL() bool {
	if u.options.BaseShortURL == "" {
		return false
	}

	parsedURL, err := url.Parse(u.options.BaseShortURL)
	if err != nil {
		return false
	}

	return parsedURL.IsAbs()
}

func (u *linkWeb) shortURL(shortCode string) string {
	return strings.TrimSuffix(u.options.BaseShortURL, "/") + "/" + shortCode
}

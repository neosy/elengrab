package link

import (
	"net/url"
	"strings"
	"uuid"

	dlink "github.com/neosy/elengrab/internal/domain/link"
)

func (u *link) GenerateShortCode(url string, length uint8, deterministic bool) string {
	return dlink.GenerateShortCode(uuid.Nil(), url, length, true)
}

func (u *link) shortURL(shortCode string) string {
	if u.options.BaseURL == "" {
		return ""
	}

	baseURL := strings.TrimSuffix(u.options.BaseURL, "/")

	return baseURL + "/" + shortCode
}

func (u *link) isAbsBaseURL() bool {
	if u.options.BaseURL == "" {
		return false
	}

	parsedURL, err := url.Parse(u.options.BaseURL)
	if err != nil {
		return false
	}

	return parsedURL.IsAbs()
}

func (u *link) shouldMatchShortURL(link *dlink.Link) bool {
	if link == nil || !link.IsMatchShortURL {
		return false
	}

	return u.isAbsBaseURL()
}

package helper

import (
	"net/url"
	"strings"
)

func isYouTube(raw string) bool {
	u, err := url.Parse(raw)
	if err != nil {
		return false
	}

	host := strings.ToLower(u.Hostname())

	switch host {
	case "youtube.com", "www.youtube.com", "m.youtube.com",
		"youtu.be", "www.youtu.be", "music.youtube.com":
		return true
	}

	// subdomains like foo.youtube.com
	return strings.HasSuffix(host, ".youtube.com")
}

package hostdetect

import (
	"net/url"
	"strings"

	"golang.org/x/net/publicsuffix"
)

func ExtractPlatformName(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}

	domain, err := publicsuffix.EffectiveTLDPlusOne(u.Hostname())
	if err != nil {
		return ""
	}

	return strings.Split(domain, ".")[0]
}

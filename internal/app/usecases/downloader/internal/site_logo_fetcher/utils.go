package logofetcher

import (
	"net/url"
	"strconv"
	"strings"
)

func parseSize(s string) int {
	// "180x180" → 180
	if s == "" {
		return 0
	}
	parts := strings.Split(s, "x")
	if len(parts) != 2 {
		return 0
	}
	v, _ := strconv.Atoi(parts[0])
	return v
}

func resolveURL(base, ref string) string {
	u, err := url.Parse(ref)
	if err != nil {
		return ref
	}
	if u.IsAbs() {
		return u.String()
	}
	b, _ := url.Parse(base)
	return b.ResolveReference(u).String()
}

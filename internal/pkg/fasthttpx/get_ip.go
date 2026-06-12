package fasthttpx

import (
	"strings"

	"github.com/valyala/fasthttp"
)

type (
	ClintIPOptions struct {
		// List of trusted proxies (your local proxy and Nginx)
		// "192.168.10.33": {}, // local proxy
		// add more internal proxies if needed
		TrustedProxies map[string]struct{}
	}

	// ClintIPOption defines a function type for configuring ClintIPOptions
	ClintIPOption func(opts *ClintIPOptions)
)

// ClintIPOptionTrustedProxies sets a l
// List of trusted proxies (your local proxy and Nginx)
func ClintIPOptionTrustedProxies(proxies []string) ClintIPOption {
	return func(opts *ClintIPOptions) {
		for _, proxy := range proxies {
			opts.TrustedProxies[proxy] = struct{}{}
		}
	}
}

// getClientIP returns the real client IP, taking into account trusted proxies
func GetClientIP(ctx *fasthttp.RequestCtx, opts ...ClintIPOption) string {
	ipOpts := ClintIPOptions{}
	for _, opt := range opts {
		opt(&ipOpts)
	}

	xff := string(ctx.Request.Header.Peek("X-Forwarded-For"))
	if xff != "" {
		ips := strings.Split(xff, ",")
		// iterate from first to last, return first IP that is NOT trusted proxy
		for _, ip := range ips {
			ip = strings.TrimSpace(ip)
			if _, ok := ipOpts.TrustedProxies[ip]; !ok {
				return ip
			}
		}
		// if all IPs are trusted proxies, return the last one
		return strings.TrimSpace(ips[len(ips)-1])
	}

	// fallback to X-Real-IP
	xRealIP := string(ctx.Request.Header.Peek("X-Real-IP"))
	if xRealIP != "" {
		return strings.TrimSpace(xRealIP)
	}

	// final fallback to RemoteAddr
	addr := ctx.RemoteAddr().String()
	ip := strings.Split(addr, ":")[0]
	return ip
}

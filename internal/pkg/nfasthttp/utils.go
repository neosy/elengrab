package nfasthttp

import (
	"strings"

	"github.com/valyala/fasthttp"
)

// List of trusted proxies (your local proxy and Nginx)
var trustedProxies = map[string]struct{}{
	// "192.168.10.33": {}, // local proxy
	// add more internal proxies if needed
}

// getClientIP returns the real client IP, taking into account trusted proxies
func getClientIP(ctx *fasthttp.RequestCtx) string {
	xff := string(ctx.Request.Header.Peek("X-Forwarded-For"))
	if xff != "" {
		ips := strings.Split(xff, ",")
		// iterate from first to last, return first IP that is NOT trusted proxy
		for _, ip := range ips {
			ip = strings.TrimSpace(ip)
			if _, ok := trustedProxies[ip]; !ok {
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

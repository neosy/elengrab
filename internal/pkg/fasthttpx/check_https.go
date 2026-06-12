package fasthttpx

import (
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/valyala/fasthttp"
)

// IsForwardedHTTPS checks if the request was made over HTTPS by looking at the "X-Forwarded-Proto" header.
// This is useful when the application is behind a reverse proxy that terminates SSL/TLS.
func IsForwardedHTTPS(ctx *fasthttp.RequestCtx) bool {
	return string(ctx.Request.Header.Peek("X-Forwarded-Proto")) == "https"
}

// EnforceHTTPS checks if the request was made over HTTPS using the IsForwardedHTTPS function.
// If the request is not over HTTPS, it returns an error indicating that HTTPS is required.
func EnforceHTTPS(ctx *fasthttp.RequestCtx) error {
	if !IsForwardedHTTPS(ctx) {
		return errorx.NewHTTPMessage("HTTPS is required", fasthttp.StatusUpgradeRequired)
	}
	return nil
}

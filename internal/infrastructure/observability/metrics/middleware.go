package metrics

import (
	"strconv"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
)

func MiddlewareHandler(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		start := time.Now().UTC()

		// save HTTP path BEFORE handler is executed
		rawPath := string(append([]byte(nil), ctx.Path()...))
		method := string(ctx.Method())

		next(ctx)

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(ctx.Response.StatusCode())

		path := normalizePath(rawPath)

		httpRequestsTotal.
			WithLabelValues(method, path, status).
			Inc()

		httpRequestDuration.
			WithLabelValues(method, path).
			Observe(duration)
	}
}

func normalizePath(p string) string {
	// /downloader/file/{uuid}/logo
	if strings.HasPrefix(p, "/downloader/file/") {
		parts := strings.Split(p, "/")
		if len(parts) >= 5 {
			parts[3] = ":id"
			return strings.Join(parts, "/")
		}
	}

	// /downloader/files/{uuid}/stream
	if strings.HasPrefix(p, "/downloader/files/") {
		parts := strings.Split(p, "/")
		if len(parts) >= 5 {
			parts[3] = ":id"
			return strings.Join(parts, "/")
		}
	}

	if strings.HasPrefix(p, "/downloader/stream/") {
		parts := strings.Split(p, "/")
		if len(parts) >= 4 {
			parts[3] = ":id"
			return strings.Join(parts, "/")
		}
	}

	if strings.HasPrefix(p, "/downloader/channel/") {
		parts := strings.Split(p, "/")
		if len(parts) >= 5 {
			parts[3] = ":id"
			return strings.Join(parts, "/")
		}
	}

	if strings.HasPrefix(p, "/s/") {
		return "/s/:id"
	}

	if strings.HasPrefix(p, "/static/") {
		return "/static/*"
	}

	return p
}

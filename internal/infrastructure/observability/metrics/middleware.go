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
	switch {
	// /admin/users/{userId}/detail
	// /admin/users/{userId}/table-row
	// /admin/users/{userId}/roles
	case strings.HasPrefix(p, "/admin/users/"):
		return normalizeBySegments(p, 3, "/admin/users/:id")

	// /downloader/items/{id}/stream
	case strings.HasPrefix(p, "/downloader/items/"):
		return normalizeBySegments(p, 3, "/downloader/items/:id")

	// /downloader/stream/{id}
	case strings.HasPrefix(p, "/downloader/stream/"):
		return normalizeBySegments(p, 3, "/downloader/stream/:id")

	// /downloader/channel/{id}/avatar
	case strings.HasPrefix(p, "/downloader/channel/"):
		return normalizeBySegments(p, 3, "/downloader/channel/:id")

	// short links
	case strings.HasPrefix(p, "/s/"):
		return "/s/:id"

	// static assets
	case strings.HasPrefix(p, "/static/"):
		return "/static/*"

	default:
		return p
	}
}

func normalizeBySegments(p string, idIndex int, fallback string) string {
	parts := strings.Split(p, "/")

	// /a/b/c => ["", "a", "b", "c"]
	if len(parts) <= idIndex {
		return fallback
	}

	parts[idIndex] = ":id"
	return strings.Join(parts, "/")
}

package fasthttpx

import "github.com/valyala/fasthttp"

type (
	Middleware        func(fasthttp.RequestHandler) fasthttp.RequestHandler
	MiddlewareWrapper func(Middleware) Middleware
)

// MiddlewareChain composes multiple middleware functions into a single middleware.
// The resulting middleware applies the given middlewares in order, wrapping the handler
// from last to first, so that the first middleware in the list is executed first
// during request handling.
//
// Example:
//
//	chain := MiddlewareChain(m1, m2, m3)
//	finalHandler := chain(h)
//
// Execution order:
//
//	m1 -> m2 -> m3 -> h
func MiddlewareChain(middlewares ...Middleware) Middleware {
	return func(h fasthttp.RequestHandler) fasthttp.RequestHandler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			h = middlewares[i](h)
		}
		return h
	}
}

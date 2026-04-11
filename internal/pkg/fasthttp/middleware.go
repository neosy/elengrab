package nfasthttp

import (
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

// Group represents a group of routes with a common prefix and shared middleware.
type RouterGroup struct {
	prefix      string
	middlewares []Middleware
	router      *router.Router
}

// NewRouterGroup creates a new Group with the given prefix and router.
// The prefix is prepended to all routes registered in this group.
// Example:
//
//	group := NewRouterGroup("/api", r)
//	group.Use(mw1, mw2)
//	group.GET("/users", usersHandler)
//
// This will register a GET route at "/api/users" with the middlewares mw1 and mw2 applied.
func NewRouterGroup(prefix string, r *router.Router) *RouterGroup {
	return &RouterGroup{
		prefix: prefix,
		router: r,
	}
}

// Use adds middleware to the group. The middlewares will be applied to all routes registered in this group.
// Example:
//
//	group := NewRouterGroup("/api", r)
//	group.Use(mw1, mw2)
//	group.GET("/users", usersHandler)
//
// This will register a GET route at "/api/users" with the middlewares mw1 and mw2 applied.
func (g *RouterGroup) Use(mw ...Middleware) {
	g.middlewares = append(g.middlewares, mw...)
}

// wrap applies the group's middleware to the given handler.
func (g *RouterGroup) wrap(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return MiddlewareChain(g.middlewares...)(h)
}

// handle registers a route with the given method, path, and handler, applying the group's prefix and middleware.
func (g *RouterGroup) handle(method, path string, h fasthttp.RequestHandler) {
	g.router.Handle(method, g.prefix+path, g.wrap(h))
}

func (g *RouterGroup) GET(path string, h fasthttp.RequestHandler) {
	g.handle(fasthttp.MethodGet, path, h)
}

func (g *RouterGroup) HEAD(path string, h fasthttp.RequestHandler) {
	g.handle(fasthttp.MethodHead, path, h)
}

func (g *RouterGroup) POST(path string, h fasthttp.RequestHandler) {
	g.handle(fasthttp.MethodPost, path, h)
}

func (g *RouterGroup) PUT(path string, h fasthttp.RequestHandler) {
	g.handle(fasthttp.MethodPut, path, h)
}

func (g *RouterGroup) PATCH(path string, h fasthttp.RequestHandler) {
	g.handle(fasthttp.MethodPatch, path, h)
}

func (g *RouterGroup) DELETE(path string, h fasthttp.RequestHandler) {
	g.handle(fasthttp.MethodDelete, path, h)
}

func (g *RouterGroup) CONNECT(path string, h fasthttp.RequestHandler) {
	g.handle(fasthttp.MethodConnect, path, h)
}

func (g *RouterGroup) OPTIONS(path string, h fasthttp.RequestHandler) {
	g.handle(fasthttp.MethodOptions, path, h)
}

func (g *RouterGroup) TRACE(path string, h fasthttp.RequestHandler) {
	g.handle(fasthttp.MethodTrace, path, h)
}

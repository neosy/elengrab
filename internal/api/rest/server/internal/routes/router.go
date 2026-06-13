package routes

import (
	"github.com/fasthttp/router"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/middleware"
)

func newRouter(middlewares *middleware.Middlewares) *router.Router {
	r := router.New()

	r.RedirectTrailingSlash = false
	r.NotFound = middlewares.Error.ErrorNotFoundHandler
	r.MethodNotAllowed = middlewares.Error.ErrorMethodNotAllowedHandler

	return r
}

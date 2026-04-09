package httpsrv

import (
	"github.com/fasthttp/router"
)

// newRouter returns a new router.
func (s *httpServer) newRouter() *router.Router {
	r := router.New()

	r.RedirectTrailingSlash = false
	r.NotFound = s.errorMiddleware.ErrorNotFoundHandler
	r.MethodNotAllowed = s.errorMiddleware.ErrorMethodNotAllowedHandler

	s.setupStaticRoutes(r, s.handlers.Static)
	s.setupUIRootRoutes(r, s.handlers.UI)
	s.setupUIRoutes(r, s.handlers.UI)
	s.setupAPIV1Routes(r, s.handlers.API)

	return r
}

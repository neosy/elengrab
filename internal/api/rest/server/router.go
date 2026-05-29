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
	s.setupUIRootRoutes(r, s.handlers.UI.Downloader)
	s.setupUIAdminRoutes(r, s.handlers.UI.Admin)
	s.setupUIDownloaderRoutes(r, s.handlers.UI.Downloader)
	s.setupAPIV1Routes(r, s.handlers.API)

	return r
}

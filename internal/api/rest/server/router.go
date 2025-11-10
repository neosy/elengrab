package httpsrv

import (
	"github.com/fasthttp/router"
	"github.com/neosy/elengrab/internal/api/rest/server/handlers"
)

// newRouter returns a new router.
func (s *httpServer) newRouter() *router.Router {
	r := router.New()

	r.RedirectTrailingSlash = false

	handlers := handlers.New(&handlers.Dependencies{
		Usecases:  s.usecases,
		AssetsDir: s.assetsDir,
	})

	s.setupHtmxRootRoutes(r, handlers.HTMX)
	s.setupHtmxUIRoutes(r, handlers.HTMX)

	return r
}

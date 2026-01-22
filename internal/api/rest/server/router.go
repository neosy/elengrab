package httpsrv

import (
	"github.com/fasthttp/router"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers"
)

// newRouter returns a new router.
func (s *httpServer) newRouter() *router.Router {
	r := router.New()

	r.RedirectTrailingSlash = false

	handlers := handlers.New(&handlers.Dependencies{
		Usecases:     s.usecases,
		Templates:    s.templates,
		AssetsDir:    s.assetsDir,
		DownloadsDir: s.downloadsDir,
	})

	s.setupStaticRoutes(r, handlers.Static)

	s.setupUIRootRoutes(r, handlers.UI)
	s.setupUIRoutes(r, handlers.UI)

	s.setupAPIV1Routes(r, handlers.API)

	return r
}

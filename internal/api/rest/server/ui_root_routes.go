package httpsrv

import (
	"github.com/fasthttp/router"
	handlers "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader_handlers"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
)

// setupUIRootRoutes setup root routes.
func (s *httpServer) setupUIRootRoutes(r *router.Router, handlers *handlers.DownloaderHandlers) {
	middlewareError := s.errorMiddleware.ErrorHandler

	// With middleware (error, auth or anonym)
	g := nfasthttp.NewRouterGroup("", r)
	g.Use(middlewareError, s.authMiddleware.AuthOrAnonym)
	{
		// Index
		g.GET(httppaths.PathIndex, handlers.IndexPageHandler)
		r.HEAD(httppaths.PathIndex, handlers.IndexPageHandler)

		// /favicon.ico
		g.GET(httppaths.PathRootFaviconICO, handlers.AssetFabiconHandler)
		r.HEAD(httppaths.PathRootFaviconICO, handlers.AssetFabiconHandler)

		// /robots.txt
		g.GET(httppaths.PathRootRobotsTxt, handlers.AssetRobotsHandler)
		g.HEAD(httppaths.PathRootRobotsTxt, handlers.AssetRobotsHandler)
	}
}

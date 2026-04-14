package httpsrv

import (
	"github.com/fasthttp/router"
	uih "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttp"
)

// setupUIRootRoutes setup root routes.
func (s *httpServer) setupUIRootRoutes(r *router.Router, handlers *uih.UIHandlers) {
	middlewareError := s.errorMiddleware.ErrorHandler

	// With middleware (error, auth or anonym)
	g := nfasthttp.NewRouterGroup("", r)
	g.Use(middlewareError, s.authMiddleware.AuthOrAnonym)
	{
		// Index
		g.GET(httppaths.PathIndex, handlers.Downloader.IndexHandler)
		r.HEAD(httppaths.PathIndex, handlers.Downloader.IndexHandler)

		// /favicon.ico
		g.GET(httppaths.PathRootFaviconICO, handlers.Downloader.RooFilesHandler)
		r.HEAD(httppaths.PathRootFaviconICO, handlers.Downloader.RooFilesHandler)

		// /robots.txt
		g.GET(httppaths.PathRootRobotsTxt, handlers.Downloader.RooFilesHandler)
		g.HEAD(httppaths.PathRootRobotsTxt, handlers.Downloader.RooFilesHandler)
	}
}

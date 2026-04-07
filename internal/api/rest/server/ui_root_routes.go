package httpsrv

import (
	"github.com/fasthttp/router"
	uih "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
)

// setupUIRootRoutes setup root routes.
func (s *httpServer) setupUIRootRoutes(r *router.Router, handlers *uih.UIHandlers) {
	authOrAnonym := s.authMiddleware.AuthOrAnonym

	// Index
	r.GET(httppaths.PathIndex, authOrAnonym(handlers.Downloader.IndexHandler))
	r.HEAD(httppaths.PathIndex, handlers.Downloader.IndexHandler)

	// /favicon.ico
	r.GET(httppaths.PathRootFaviconICO, authOrAnonym(handlers.Downloader.RooFilesHandler))
	r.HEAD(httppaths.PathRootFaviconICO, handlers.Downloader.RooFilesHandler)

	// /robots.txt
	r.GET(httppaths.PathRootRobotsTxt, authOrAnonym(handlers.Downloader.RooFilesHandler))
	r.HEAD(httppaths.PathRootRobotsTxt, authOrAnonym(handlers.Downloader.RooFilesHandler))
}

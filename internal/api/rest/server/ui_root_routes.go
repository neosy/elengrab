package httpsrv

import (
	"github.com/fasthttp/router"
	uih "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
)

// setupUIRootRoutes setup root routes.
func (s *httpServer) setupUIRootRoutes(r *router.Router, handlers *uih.UIHandlers) {
	optAuth := s.authMiddleware.OptionalAuth

	// Index
	r.GET(httppaths.PathIndex, optAuth(handlers.Downloader.IndexHandler))

	// /favicon.ico
	r.GET(httppaths.PathRootFaviconICO, handlers.Downloader.FaviconICOHandler)
}

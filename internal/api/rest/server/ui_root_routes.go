package httpsrv

import (
	"github.com/fasthttp/router"
	statich "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/static"
	uih "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
)

// setupUIRootRoutes setup root routes.
func (s *httpServer) setupUIRootRoutes(r *router.Router, handlers *uih.UIHandlers, staticHandlers *statich.StaticHandlers) {
	auth := s.authMiddleware.AutoRegister

	// Index
	r.GET(httppaths.PathIndex, auth(handlers.Downloader.IndexHandler))

	// /favicon.ico
	r.GET(httppaths.PathRootFaviconICO, handlers.Downloader.FaviconICOHandler)
}

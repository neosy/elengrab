package httpsrv

import (
	"github.com/fasthttp/router"
	uih "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
)

// setupUIRootRoutes setup root routes.
func (s *httpServer) setupUIRootRoutes(r *router.Router, handlers *uih.UIHandlers) {
	auth := s.authMiddleware.AutoRegister

	// Index
	r.GET(httppaths.PathIndex, auth(handlers.Downloader.IndexHandler))
}

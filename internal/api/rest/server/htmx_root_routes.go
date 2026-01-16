package httpsrv

import (
	"github.com/fasthttp/router"
	htmxh "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/htmx"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
)

// setupHtmxRootoutes setup root routes.
func (s *httpServer) setupHtmxRootRoutes(r *router.Router, handlers *htmxh.HTMXHandlers) {
	auth := s.authMiddleware.AutoRegister

	// Index
	r.GET(httppaths.PathIndex, auth(handlers.Grabber.IndexHandler))
}

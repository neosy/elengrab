package httpsrv

import (
	"github.com/fasthttp/router"
	htmxh "github.com/neosy/elengrab/internal/api/rest/server/handlers/htmx"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/paths"
)

// setupHtmxRootoutes setup root routes.
func (s *httpServer) setupHtmxRootRoutes(r *router.Router, handlers *htmxh.Handlers) {
	// Index
	r.GET(httppaths.PathIndex, handlers.Grabber.IndexHandler)
}

package httpsrv

import (
	"github.com/fasthttp/router"
	statich "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/static"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
)

// setupStaticRoutes setup Static routes.
func (s *httpServer) setupStaticRoutes(r *router.Router, handlers *statich.StaticHandlers) {
	// Static
	group := r.Group(httppaths.GroupStatic)
	{
		group.GET(httppaths.PathCssFiles, handlers.Static.StaticCssHandler)
		group.GET(httppaths.PathImgFiles, handlers.Static.StaticImgHandler)
		group.GET(httppaths.PathJsFiles, handlers.Static.StaticJsHandler)
		group.GET(httppaths.PathPwaFiles, handlers.Static.StaticPwaHandler)
	}
}

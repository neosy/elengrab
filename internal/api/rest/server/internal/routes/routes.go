package routes

import (
	"github.com/fasthttp/router"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/middleware"
)

type routes struct {
	middlewares *middleware.Middlewares
	router      *router.Router
}

func NewRoutes(
	middlewares *middleware.Middlewares,
	handlers *handlers.Handlers,
	// Options
	shortLinkPrefix string,
) *routes {
	routes := &routes{
		middlewares: middlewares,
		router:      newRouter(middlewares),
	}

	routes.registerStatic(handlers.Static)
	routes.registerUIRoot(handlers.UI.Downloader)
	routes.registerUIAdmin(handlers.UI.Admin)
	routes.registerUIDownloader(handlers.UI.Downloader, shortLinkPrefix)
	routes.registerAPI(handlers.API.V1)

	return routes
}

// Router returns a router.
func (r *routes) Router() *router.Router {
	return r.router
}

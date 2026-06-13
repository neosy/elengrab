package routes

import (
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
)

// registerUIRoot register ui root routes.
func (r *routes) registerUIRoot(handlers *downloader.DownloaderHandlers) {
	middlewareError := r.middlewares.Error.ErrorHandler

	// With middleware (error, auth or anonym)
	g := nfasthttp.NewRouterGroup("", r.router)
	g.Use(middlewareError, r.middlewares.Auth.AuthOrAnonym)
	{
		// Index
		g.GET(httppaths.PathIndex, handlers.IndexPageHandler)
		r.router.HEAD(httppaths.PathIndex, handlers.IndexPageHandler)

		// /favicon.ico
		g.GET(httppaths.PathRootFaviconICO, handlers.AssetFabiconHandler)
		r.router.HEAD(httppaths.PathRootFaviconICO, handlers.AssetFabiconHandler)

		// /robots.txt
		g.GET(httppaths.PathRootRobotsTxt, handlers.AssetRobotsHandler)
		g.HEAD(httppaths.PathRootRobotsTxt, handlers.AssetRobotsHandler)
	}
}

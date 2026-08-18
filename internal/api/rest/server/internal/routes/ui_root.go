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
		g.GET(httppaths.IndexPath, handlers.IndexPageHandler)
		r.router.HEAD(httppaths.IndexPath, handlers.IndexPageHandler)

		// /favicon.ico
		g.GET(httppaths.RootFaviconICOPath, handlers.AssetFabiconHandler)
		r.router.HEAD(httppaths.RootFaviconICOPath, handlers.AssetFabiconHandler)

		// /robots.txt
		g.GET(httppaths.RootRobotsTxtPath, handlers.AssetRobotsHandler)
		g.HEAD(httppaths.RootRobotsTxtPath, handlers.AssetRobotsHandler)
	}
}

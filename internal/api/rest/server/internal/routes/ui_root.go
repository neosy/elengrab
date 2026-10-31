package routes

import (
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	"github.com/neosy/elengrab/internal/pkg/fasthttpx"
)

// registerUIRoot register ui root routes.
func (r *routes) registerUIRoot(handlers *downloader.DownloaderHandlers) {
	middlewareError := r.middlewares.Error.ErrorHandler

	// With middleware (error, auth or anonym)
	g := fasthttpx.NewRouterGroup("", r.router)
	g.Use(middlewareError, r.middlewares.Auth.AuthOrAnonym)
	{
		// Index
		r.router.HEAD(httppaths.IndexPath, fasthttpx.HTMLHeadHandler)
		g.GET(httppaths.IndexPath, handlers.IndexPageHandler)

		// /favicon.ico
		r.router.HEAD(httppaths.RootFaviconICOPath, fasthttpx.HTMLHeadHandler)
		g.GET(httppaths.RootFaviconICOPath, handlers.AssetFaviconHandler)

		// /robots.txt
		r.router.HEAD(httppaths.RootRobotsTxtPath, fasthttpx.HTMLHeadHandler)
		g.GET(httppaths.RootRobotsTxtPath, handlers.AssetRobotsHandler)
	}

	// User downloads
	{
		// With middleware (error, auth or anonym)
		g := fasthttpx.NewRouterGroup("", r.router)
		g.Use(middlewareError, r.middlewares.Auth.AuthOrAnonym)
		{
			g.HEAD(httppaths.RootUserLoginPath, fasthttpx.HTMLHeadHandler)
			g.GET(httppaths.RootUserLoginPath, handlers.UserDownloadsPageHandler)
		}
	}
}

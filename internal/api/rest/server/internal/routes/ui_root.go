package routes

import (
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	"github.com/neosy/elengrab/internal/pkg/fasthttpx"
)

// registerUIRoot register ui root routes.
func (r *routes) registerUIRoot(handlers *downloader.DownloaderHandlers) {
	middlewareError := r.middlewares.Error.ErrorHandler

	// Without middleware
	g := fasthttpx.NewRouterGroup("", r.router)
	{
		// /favicon.ico
		g.HEAD(httppaths.RootFaviconICOPath, handlers.AssetFaviconHandler)
		g.GET(httppaths.RootFaviconICOPath, handlers.AssetFaviconHandler)

		// /robots.txt
		g.HEAD(httppaths.RootRobotsTxtPath, handlers.AssetRobotsHandler)
		g.GET(httppaths.RootRobotsTxtPath, handlers.AssetRobotsHandler)
	}

	// With middleware (error, auth or anonym)
	g = fasthttpx.NewRouterGroup("", r.router)
	g.Use(middlewareError, r.middlewares.Auth.AuthOrAnonym)
	{
		// Index
		g.HEAD(httppaths.IndexPath, handlers.IndexPageHandler)
		g.GET(httppaths.IndexPath, handlers.IndexPageHandler)
	}

	// User downloads
	{
		// With middleware (error, auth or anonym)
		g := fasthttpx.NewRouterGroup("", r.router)
		g.Use(middlewareError, r.middlewares.Auth.AuthOrAnonym)
		{
			g.HEAD(httppaths.RootUserLoginPath, handlers.UserDownloadsPageHandler)
			g.GET(httppaths.RootUserLoginPath, handlers.UserDownloadsPageHandler)
		}
	}
}

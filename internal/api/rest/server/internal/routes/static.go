package routes

import (
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/static"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
)

// registerStatic register Static routes.
func (r *routes) registerStatic(handlers *static.StaticHandlers) {
	middlewareError := r.middlewares.Error.ErrorHandler

	// Static with middleware (error)
	g := nfasthttp.NewRouterGroup(httppaths.StaticGroup, r.router)
	g.Use(middlewareError)
	{
		g.GET(httppaths.CssFilesPath, handlers.StaticCssHandler)
		g.HEAD(httppaths.CssFilesPath, handlers.StaticCssHandler)
		g.GET(httppaths.FontFilesPath, handlers.StaticFontsHandler)
		g.HEAD(httppaths.FontFilesPath, handlers.StaticFontsHandler)
		g.GET(httppaths.ImageFilesPath, handlers.StaticImgHandler)
		g.HEAD(httppaths.ImageFilesPath, handlers.StaticImgHandler)
		g.GET(httppaths.IconFilesPath, handlers.StaticIconHandler)
		g.HEAD(httppaths.IconFilesPath, handlers.StaticIconHandler)
		g.GET(httppaths.JsFilesPath, handlers.StaticJsHandler)
		g.HEAD(httppaths.JsFilesPath, handlers.StaticJsHandler)
		g.GET(httppaths.PwaFilesPath, handlers.StaticPwaHandler)
		g.HEAD(httppaths.PwaFilesPath, handlers.StaticPwaHandler)
		g.GET(httppaths.ThumbnailPath, handlers.ThumbnailHandler)
		g.HEAD(httppaths.ThumbnailPath, handlers.ThumbnailHandler)
		g.GET(httppaths.YoutubeChannelPath, handlers.YoutubeChannelHandler)
		g.HEAD(httppaths.YoutubeChannelPath, handlers.YoutubeChannelHandler)
	}
}

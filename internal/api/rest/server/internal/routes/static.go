package routes

import (
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/static"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	"github.com/neosy/elengrab/internal/pkg/fasthttpx"
)

// registerStatic register Static routes.
func (r *routes) registerStatic(handlers *static.StaticHandlers) {
	middlewareError := r.middlewares.Error.ErrorHandler

	// Static with middleware (error)
	g := fasthttpx.NewRouterGroup(httppaths.StaticGroup, r.router)
	g.Use(middlewareError)
	{
		g.HEAD(httppaths.CssFilesPath, fasthttpx.CSSHeadHandler)
		g.GET(httppaths.CssFilesPath, handlers.StaticCssHandler)

		g.HEAD(httppaths.JsFilesPath, fasthttpx.JSHeadHandler)
		g.GET(httppaths.JsFilesPath, handlers.StaticJsHandler)

		g.HEAD(httppaths.FontFilesPath, fasthttpx.HTMLHeadHandler)
		g.GET(httppaths.FontFilesPath, handlers.StaticFontsHandler)

		g.HEAD(httppaths.PwaFilesPath, fasthttpx.HTMLHeadHandler)
		g.GET(httppaths.PwaFilesPath, handlers.StaticPwaHandler)

		g.HEAD(httppaths.ImageFilesPath, fasthttpx.HTMLHeadHandler)
		g.GET(httppaths.ImageFilesPath, handlers.StaticImgHandler)

		g.HEAD(httppaths.IconFilesPath, fasthttpx.HTMLHeadHandler)
		g.GET(httppaths.IconFilesPath, handlers.StaticIconHandler)

		g.HEAD(httppaths.ThumbnailPath, fasthttpx.HTMLHeadHandler)
		g.GET(httppaths.ThumbnailPath, handlers.ThumbnailHandler)

		g.HEAD(httppaths.ChannelImagePath, fasthttpx.HTMLHeadHandler)
		g.GET(httppaths.ChannelImagePath, handlers.ChannelImageHandler)
	}
}

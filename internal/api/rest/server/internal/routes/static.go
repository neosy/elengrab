package routes

import (
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/static"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	"github.com/neosy/elengrab/internal/pkg/fasthttpx"
)

// registerStatic register Static routes.
func (r *routes) registerStatic(handlers *static.StaticHandlers) {
	// Static with middleware (error)
	g := fasthttpx.NewRouterGroup(httppaths.StaticGroup, r.router)
	{
		g.HEAD(httppaths.CssFilesPath, fasthttpx.CSSHeadHandler)
		g.GET(httppaths.CssFilesPath, handlers.StaticCssHandler)

		g.HEAD(httppaths.JsFilesPath, fasthttpx.JSHeadHandler)
		g.GET(httppaths.JsFilesPath, handlers.StaticJsHandler)

		g.HEAD(httppaths.FontFilesPath, handlers.StaticFontsHandler)
		g.GET(httppaths.FontFilesPath, handlers.StaticFontsHandler)

		g.HEAD(httppaths.PwaFilesPath, handlers.StaticPwaHandler)
		g.GET(httppaths.PwaFilesPath, handlers.StaticPwaHandler)

		g.HEAD(httppaths.ImageFilesPath, handlers.StaticImgHandler)
		g.GET(httppaths.ImageFilesPath, handlers.StaticImgHandler)

		g.HEAD(httppaths.IconFilesPath, handlers.StaticIconHandler)
		g.GET(httppaths.IconFilesPath, handlers.StaticIconHandler)

		g.HEAD(httppaths.ThumbnailPath, handlers.ThumbnailHandler)
		g.GET(httppaths.ThumbnailPath, handlers.ThumbnailHandler)

		g.HEAD(httppaths.ChannelImagePath, handlers.ChannelImageHandler)
		g.GET(httppaths.ChannelImagePath, handlers.ChannelImageHandler)
	}
}

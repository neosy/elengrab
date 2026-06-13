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
	g := nfasthttp.NewRouterGroup(httppaths.GroupStatic, r.router)
	g.Use(middlewareError)
	{
		g.GET(httppaths.PathCssFiles, handlers.StaticCssHandler)
		g.HEAD(httppaths.PathCssFiles, handlers.StaticCssHandler)
		g.GET(httppaths.PathFontFiles, handlers.StaticFontsHandler)
		g.HEAD(httppaths.PathFontFiles, handlers.StaticFontsHandler)
		g.GET(httppaths.PathImageFiles, handlers.StaticImgHandler)
		g.HEAD(httppaths.PathImageFiles, handlers.StaticImgHandler)
		g.GET(httppaths.PathIconFiles, handlers.StaticIconHandler)
		g.HEAD(httppaths.PathIconFiles, handlers.StaticIconHandler)
		g.GET(httppaths.PathJsFiles, handlers.StaticJsHandler)
		g.HEAD(httppaths.PathJsFiles, handlers.StaticJsHandler)
		g.GET(httppaths.PathPwaFiles, handlers.StaticPwaHandler)
		g.HEAD(httppaths.PathPwaFiles, handlers.StaticPwaHandler)
		g.GET(httppaths.PathThumbnail, handlers.ThumbnailHandler)
		g.HEAD(httppaths.PathThumbnail, handlers.ThumbnailHandler)
		g.GET(httppaths.PathYoutubeChannel, handlers.YoutubeChannelHandler)
		g.HEAD(httppaths.PathYoutubeChannel, handlers.YoutubeChannelHandler)
	}
}

package httpsrv

import (
	"github.com/fasthttp/router"
	statich "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/static"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttp"
)

// setupStaticRoutes setup Static routes.
func (s *httpServer) setupStaticRoutes(r *router.Router, handlers *statich.StaticHandlers) {
	middlewareError := s.errorMiddleware.ErrorHandler

	// Static with middleware (error)
	g := nfasthttp.NewRouterGroup(httppaths.GroupStatic, r)
	g.Use(middlewareError)
	{
		g.GET(httppaths.PathCssFiles, handlers.Static.StaticCssHandler)
		g.GET(httppaths.PathImgFiles, handlers.Static.StaticImgHandler)
		g.GET(httppaths.PathIconFiles, handlers.Static.StaticIconHandler)
		g.GET(httppaths.PathJsFiles, handlers.Static.StaticJsHandler)
		g.GET(httppaths.PathPwaFiles, handlers.Static.StaticPwaHandler)
		g.GET(httppaths.PathThumbnail, handlers.Static.ThumbnailHandler)
		g.GET(httppaths.PathYoutubeChannel, handlers.Static.YoutubeChannelHandler)
	}

	// Static without middleware
	gg := r.Group(httppaths.GroupStatic)
	{
		gg.HEAD(httppaths.PathCssFiles, handlers.Static.StaticCssHandler)
		gg.HEAD(httppaths.PathImgFiles, handlers.Static.StaticImgHandler)
		gg.HEAD(httppaths.PathIconFiles, handlers.Static.StaticIconHandler)
		gg.HEAD(httppaths.PathJsFiles, handlers.Static.StaticJsHandler)
		gg.HEAD(httppaths.PathPwaFiles, handlers.Static.StaticPwaHandler)
		gg.HEAD(httppaths.PathThumbnail, handlers.Static.ThumbnailHandler)
		gg.HEAD(httppaths.PathYoutubeChannel, handlers.Static.YoutubeChannelHandler)
	}

}

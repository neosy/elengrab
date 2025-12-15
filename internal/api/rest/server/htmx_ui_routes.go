package httpsrv

import (
	"github.com/fasthttp/router"
	htmxh "github.com/neosy/elengrab/internal/api/rest/server/handlers/htmx"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/paths"
)

// setupHtmxUIRoutes setup UI routes.
func (s *httpServer) setupHtmxUIRoutes(r *router.Router, handlers *htmxh.Handlers) {
	// Static
	group := r.Group(httppaths.GroupStatic)
	{
		group.GET(httppaths.PathCssFiles, handlers.Static.StaticCssHandler)
		group.GET(httppaths.PathImgFiles, handlers.Static.StaticImgHandler)
		group.GET(httppaths.PathJsFiles, handlers.Static.StaticJsHandler)
		group.GET(httppaths.PathPwaFiles, handlers.Static.StaticPwaHandler)
	}

	// Downloader
	group = r.Group(httppaths.GroupDownloader)
	{
		group.GET(httppaths.PathHistory, handlers.Grabber.GetFilesHistoryHandler)
		group.POST(httppaths.PathGrab, handlers.Grabber.GrabHandler)
		group.GET(httppaths.PathDownload, handlers.Grabber.DownloadHandler)
		group.GET(httppaths.PathFileRow, handlers.Grabber.GetFileRow)
		group.DELETE(httppaths.PathFileRow, handlers.Grabber.DeleteFileRow)
		group.GET(httppaths.PathChannelAvatar, handlers.Grabber.GetChannelAvatar)
	}
}

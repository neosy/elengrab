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
		group.GET(httppaths.PathFiles, handlers.Static.StaticHandler)
	}

	// Downloader
	group = r.Group(httppaths.GroupDownloader)
	{
		group.POST(httppaths.PathGrab, handlers.Grabber.GrabHandler)
		group.GET(httppaths.PathDownload, handlers.Grabber.DownloadHandler)
	}
}

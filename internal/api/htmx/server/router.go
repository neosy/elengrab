package httpsrv

import (
	"github.com/fasthttp/router"
	"github.com/neosy/elengrab/internal/api/htmx/server/handlers"
)

const (
	// Groups
	groupStatic     = "/static"
	groupDownloader = "/downloader"

	pathIndex = "/"
	// Counter increment
	pathGrab     = "/grab"
	pathDownload = "/download"
)

// newRouter returns a new router.
func (s *httpServer) newRouter(assetsDir string) *router.Router {
	r := router.New()

	r.RedirectTrailingSlash = false

	handlers := handlers.New(&handlers.Dependencies{
		Usecases:  s.usecases,
		AssetsDir: s.assetsDir,
	})

	r.GET(pathIndex, handlers.Index.IndexHandler)

	group := r.Group(groupStatic)
	{
		group.GET("/{filepath:*}", handlers.Static.StaticHandler)
	}

	group = r.Group(groupDownloader)
	{
		group.POST(pathGrab, handlers.Grab.GrabHandler)
		group.GET(pathDownload, handlers.Grab.DownloadHandler)
	}

	return r
}

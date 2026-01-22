package httpsrv

import (
	"github.com/fasthttp/router"
	uih "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
)

// setupUIRoutes setup UI routes.
func (s *httpServer) setupUIRoutes(r *router.Router, handlers *uih.UIHandlers) {
	auth := s.authMiddleware.AutoRegister

	// Downloader
	group := r.Group(httppaths.GroupDownloader)
	{
		group.GET(httppaths.PathHistory, auth(handlers.Downloader.GetFilesHistoryHandler))
		group.POST(httppaths.PathGrab, auth(handlers.Downloader.GrabHandler))
		group.GET(httppaths.PathDownload, auth(handlers.Downloader.DownloadHandler))
		group.GET(httppaths.PathFileRow, auth(handlers.Downloader.GetFileRow))
		group.DELETE(httppaths.PathFileRow, auth(handlers.Downloader.DeleteFileRow))
		group.GET(httppaths.PathChannelAvatar, handlers.Downloader.GetChannelAvatar)
	}
}

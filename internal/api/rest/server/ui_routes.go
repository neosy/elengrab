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
		group.GET(httppaths.PathStream, auth(handlers.Downloader.StreamHandler))
		group.GET(httppaths.PathFileRow, auth(handlers.Downloader.GetFileRowHandler))
		group.GET(httppaths.PathFileLogo, auth(handlers.Downloader.GetFileLogoHandler))
		group.DELETE(httppaths.PathFile, auth(handlers.Downloader.DeleteFileRowHandler))
		group.POST(httppaths.PathFileDownloadRepeat, auth(handlers.Downloader.RepeatDownloadHandler))
		group.GET(httppaths.PathChannelAvatar, handlers.Downloader.GetChannelAvatarHandler)
		group.GET(httppaths.PathFilesEvents, auth(handlers.Downloader.EventsHandler))
	}
}

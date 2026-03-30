package httpsrv

import (
	"github.com/fasthttp/router"
	uih "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
)

// setupUIRoutes setup UI routes.
func (s *httpServer) setupUIRoutes(r *router.Router, handlers *uih.UIHandlers) {
	reqAuth := s.authMiddleware.RequireAuth
	optAuth := s.authMiddleware.OptionalAuth

	// Account
	group := r.Group(httppaths.GroupAccount)
	{
		group.GET(httppaths.PathRegister, optAuth(handlers.Downloader.AuthRegisterHandler))
		group.POST(httppaths.PathRegister, optAuth(handlers.Downloader.AuthRegisterSubmitHandler))
		group.GET(httppaths.PathLogin, optAuth(handlers.Downloader.AuthLoginHandler))
		group.POST(httppaths.PathLogin, optAuth(handlers.Downloader.AuthLoginSubmitHandler))
		group.GET(httppaths.PathLogout, optAuth(handlers.Downloader.AuthLogoutHandler))
	}

	// Downloader
	group = r.Group(httppaths.GroupDownloader)
	{
		group.GET(httppaths.PathHistory, optAuth(handlers.Downloader.GetFilesHistoryHandler))
		group.POST(httppaths.PathGrab, reqAuth(handlers.Downloader.GrabHandler))
		group.GET(httppaths.PathDownload, optAuth(handlers.Downloader.DownloadHandler))
		group.GET(httppaths.PathStream, optAuth(handlers.Downloader.StreamHandler))
		group.GET(httppaths.PathFileRow, optAuth(handlers.Downloader.GetFileRowHandler))
		group.GET(httppaths.PathFileLogo, optAuth(handlers.Downloader.GetFileLogoHandler))
		group.DELETE(httppaths.PathFile, reqAuth(handlers.Downloader.DeleteFileRowHandler))
		group.POST(httppaths.PathFileDownloadRepeat, reqAuth(handlers.Downloader.RepeatDownloadHandler))
		group.GET(httppaths.PathChannelAvatar, handlers.Downloader.GetChannelAvatarHandler)
		group.GET(httppaths.PathFilesEvents, optAuth(handlers.Downloader.EventsHandler))
		group.POST(httppaths.PathSearch, optAuth(handlers.Downloader.SearchHandler))
	}
}

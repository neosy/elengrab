package httpsrv

import (
	"github.com/fasthttp/router"
	uih "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
)

// setupUIRoutes setup UI routes.
func (s *httpServer) setupUIRoutes(r *router.Router, handlers *uih.UIHandlers) {
	authOrGuest := s.authMiddleware.AuthOrGuest
	authOrAnonym := s.authMiddleware.AuthOrAnonym
	authOptional := s.authMiddleware.AuthOptional
	requireAuth := s.authMiddleware.RequireAuth
	requireAuthMode := s.authMiddleware.RequireAuthMode

	// Account
	group := r.Group(httppaths.GroupAccount)
	{
		group.GET(httppaths.PathRegister, authOrAnonym(handlers.Downloader.AuthRegisterHandler))
		group.POST(httppaths.PathRegister, authOrAnonym(handlers.Downloader.AuthRegisterSubmitHandler))
		group.GET(httppaths.PathLogin, authOrAnonym(handlers.Downloader.AuthLoginHandler))
		group.POST(httppaths.PathLogin, authOrAnonym(handlers.Downloader.AuthLoginSubmitHandler))
		group.GET(httppaths.PathLogout, authOrAnonym(handlers.Downloader.AuthLogoutHandler))
	}

	// Downloader
	group = r.Group(httppaths.GroupDownloader)
	{
		group.GET(httppaths.PathAccountMenu, requireAuth(handlers.Downloader.AccountMenuHandler))
		group.GET(httppaths.PathHistory, authOptional(handlers.Downloader.GetFilesHistoryHandler))
		group.POST(httppaths.PathGrab, authOrGuest(handlers.Downloader.GrabHandler))
		group.GET(httppaths.PathDownload, requireAuthMode(handlers.Downloader.DownloadHandler))
		group.GET(httppaths.PathStream, requireAuthMode(handlers.Downloader.StreamHandler))
		group.GET(httppaths.PathStreamShortCode, requireAuthMode(handlers.Downloader.StreamShortCodeHandler))
		group.GET(httppaths.PathFileRow, requireAuthMode(handlers.Downloader.GetFileRowHandler))
		group.GET(httppaths.PathFileLogo, requireAuthMode(handlers.Downloader.GetFileLogoHandler))
		group.DELETE(httppaths.PathFile, authOrGuest(handlers.Downloader.DeleteFileRowHandler))
		group.POST(httppaths.PathFileDownloadRepeat, authOrGuest(handlers.Downloader.RepeatDownloadHandler))
		group.GET(httppaths.PathChannelAvatar, authOrAnonym(handlers.Downloader.GetChannelAvatarHandler))
		group.GET(httppaths.PathFilesEvents, authOptional(handlers.Downloader.EventsHandler))
		group.POST(httppaths.PathSearch, authOptional(handlers.Downloader.SearchHandler))
		group.GET(httppaths.PathFileMenu, requireAuthMode(handlers.Downloader.RowMenuHandler))
		group.POST(httppaths.PathFileShortLink, requireAuthMode(handlers.Downloader.GetFileShortLinkHandler))
	}

	// Short link
	group = r.Group(s.shortLinkPrefix)
	{
		group.GET(httppaths.PathShortLink, handlers.Downloader.ShortLinkHandler)
	}
}

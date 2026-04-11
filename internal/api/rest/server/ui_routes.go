package httpsrv

import (
	"github.com/fasthttp/router"
	uih "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttp"
)

// setupUIRoutes setup UI routes.
func (s *httpServer) setupUIRoutes(r *router.Router, handlers *uih.UIHandlers) {
	middlewareError := s.errorMiddleware.ErrorHandler

	// Account with middleware (error, auth or anonym)
	g := nfasthttp.NewRouterGroup(httppaths.GroupAccount, r)
	g.Use(middlewareError, s.authMiddleware.AuthOrAnonym)
	{
		g.GET(httppaths.PathRegister, handlers.Downloader.AuthRegisterHandler)
		g.POST(httppaths.PathRegister, handlers.Downloader.AuthRegisterSubmitHandler)

		g.GET(httppaths.PathLogin, handlers.Downloader.AuthLoginHandler)
		g.POST(httppaths.PathLogin, handlers.Downloader.AuthLoginSubmitHandler)

		g.GET(httppaths.PathLogout, handlers.Downloader.AuthLogoutHandler)
	}

	// Downloader
	//group = r.Group(httppaths.GroupDownloader)
	{
		// With middleware (error, require auth)
		g := nfasthttp.NewRouterGroup(httppaths.GroupDownloader, r)
		g.Use(middlewareError, s.authMiddleware.RequireAuth)
		{
			g.GET(httppaths.PathAccountMenu, handlers.Downloader.AccountMenuHandler)
		}

		// With middleware (error, auth optional)
		g = nfasthttp.NewRouterGroup(httppaths.GroupDownloader, r)
		g.Use(middlewareError, s.authMiddleware.AuthOptional)
		{
			g.GET(httppaths.PathHistory, handlers.Downloader.GetFilesHistoryHandler)
			g.GET(httppaths.PathFilesEvents, handlers.Downloader.EventsHandler)
			g.POST(httppaths.PathSearch, handlers.Downloader.SearchHandler)
		}

		// With middleware (error, auth or guest)
		g = nfasthttp.NewRouterGroup(httppaths.GroupDownloader, r)
		g.Use(middlewareError, s.authMiddleware.AuthOrGuest)
		{
			g.POST(httppaths.PathGrab, handlers.Downloader.GrabHandler)
			g.DELETE(httppaths.PathFile, handlers.Downloader.DeleteFileRowHandler)
			g.POST(httppaths.PathFileDownloadRepeat, handlers.Downloader.RepeatDownloadHandler)
		}

		// With middleware (error, require auth mode)
		g = nfasthttp.NewRouterGroup(httppaths.GroupDownloader, r)
		g.Use(middlewareError, s.authMiddleware.RequireAuthMode)
		{
			g.GET(httppaths.PathFileRow, handlers.Downloader.GetFileRowHandler)
			g.GET(httppaths.PathFileLogo, handlers.Downloader.GetFileLogoHandler)
			g.GET(httppaths.PathFileMenu, handlers.Downloader.RowMenuHandler)
			g.POST(httppaths.PathFileShortLink, handlers.Downloader.GetFileShortLinkHandler)
			g.GET(httppaths.PathStream, handlers.Downloader.StreamHandler)

			g.GET(httppaths.PathDownload, handlers.Downloader.DownloadHandler)
			g.HEAD(httppaths.PathDownload, handlers.Downloader.DownloadHandler)

			g.GET(httppaths.PathStreamShortCode, handlers.Downloader.StreamShortCodeHandler)
			g.HEAD(httppaths.PathStreamShortCode, handlers.Downloader.StreamShortCodeHandler)
		}

		// With middleware (error, auth or anonym)
		g = nfasthttp.NewRouterGroup(httppaths.GroupDownloader, r)
		g.Use(middlewareError, s.authMiddleware.AuthOrAnonym)
		{
			g.GET(httppaths.PathChannelAvatar, handlers.Downloader.GetChannelAvatarHandler)
		}

	}

	// Short link
	{
		// With middleware (error, auth or anonym)
		g := nfasthttp.NewRouterGroup(s.shortLinkPrefix, r)
		g.Use(middlewareError, s.authMiddleware.AuthOrAnonym)
		{
			g.GET(httppaths.PathShortLink, handlers.Downloader.ShortLinkHandler)
		}

		// With middleware (error)
		g = nfasthttp.NewRouterGroup(s.shortLinkPrefix, r)
		g.Use(middlewareError)
		{
			g.HEAD(httppaths.PathShortLink, handlers.Downloader.ShortLinkHandler)
		}
	}
}

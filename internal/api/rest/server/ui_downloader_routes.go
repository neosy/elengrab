package httpsrv

import (
	"github.com/fasthttp/router"
	handlers "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader_handlers"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttp"
)

// setupUIDownloaderRoutes setup UI routes.
func (s *httpServer) setupUIDownloaderRoutes(r *router.Router, handlers *handlers.DownloaderHandlers) {
	middlewareError := s.errorMiddleware.ErrorHandler

	// Account with middleware (error, auth or anonym)
	g := nfasthttp.NewRouterGroup(httppaths.GroupAccount, r)
	g.Use(middlewareError, s.authMiddleware.AuthOrAnonym)
	{
		g.GET(httppaths.PathRegister, handlers.AuthRegisterPageHandler)
		g.HEAD(httppaths.PathRegister, handlers.AuthRegisterPageHandler)
		g.POST(httppaths.PathRegister, handlers.AuthRegisterSubmitHandler)

		g.GET(httppaths.PathLogin, handlers.AuthLoginPageHandler)
		g.HEAD(httppaths.PathLogin, handlers.AuthLoginPageHandler)
		g.POST(httppaths.PathLogin, handlers.AuthLoginSubmitHandler)

		g.GET(httppaths.PathLogout, handlers.AuthLogoutHandler)
	}

	// Downloader
	//group = r.Group(httppaths.GroupDownloader)
	{
		// With middleware (error, require auth)
		g := nfasthttp.NewRouterGroup(httppaths.GroupDownloader, r)
		g.Use(middlewareError, s.authMiddleware.RequireAuth)
		{
			g.GET(httppaths.PathAccountMenu, handlers.AccountMenuHandler)
			g.HEAD(httppaths.PathAccountMenu, handlers.AccountMenuHandler)
		}

		// With middleware (error, auth optional)
		g = nfasthttp.NewRouterGroup(httppaths.GroupDownloader, r)
		g.Use(middlewareError, s.authMiddleware.AuthOptional)
		{
			g.GET(httppaths.PathHistory, handlers.MediaHistoryHandler)
			g.HEAD(httppaths.PathHistory, handlers.MediaHistoryHandler)
			g.GET(httppaths.PathEvents, handlers.EventsStreamHandler)
			g.POST(httppaths.PathSearch, handlers.SearchHandler)
		}

		// With middleware (error, auth or guest)
		g = nfasthttp.NewRouterGroup(httppaths.GroupDownloader, r)
		g.Use(middlewareError, s.authMiddleware.AuthOrGuest)
		{
			g.POST(httppaths.PathGrab, handlers.ImportMediaByURLHandler)
			g.GET(httppaths.PathShareTarget, handlers.ImportFromShareHandler)
			g.DELETE(httppaths.PathMediaItem, handlers.MediaItemDeleteHandler)
			g.POST(httppaths.PathMediaItemDownloadRepeat, handlers.RetryImportMediaHandler)
		}

		// With middleware (error, require auth mode)
		g = nfasthttp.NewRouterGroup(httppaths.GroupDownloader, r)
		g.Use(middlewareError, s.authMiddleware.RequireAuthMode)
		{
			g.GET(httppaths.PathMediaItemRow, handlers.MediaItemRowHandler)
			g.HEAD(httppaths.PathMediaItemRow, handlers.MediaItemRowHandler)
			g.GET(httppaths.PathMediaItemImage, handlers.MediaItemImageHandler)
			g.HEAD(httppaths.PathMediaItemImage, handlers.MediaItemImageHandler)
			g.GET(httppaths.PathMediaItemMenu, handlers.MediaItemRowMenuHandler)
			g.POST(httppaths.PathMediaItemShortLink, handlers.CreateMediaShareLinkHandler)
			g.GET(httppaths.PathMediaItemStream, handlers.MediaItemStreamHandler)
			g.HEAD(httppaths.PathMediaItemStream, handlers.MediaItemStreamHandler)
			g.GET(httppaths.PathMediaItemWatch, handlers.WatchPageByDownloadIDHandler)
			g.HEAD(httppaths.PathMediaItemWatch, handlers.WatchPageByDownloadIDHandler)

			g.GET(httppaths.PathDownloadFile, handlers.DownloadFileHandler)
			g.HEAD(httppaths.PathDownloadFile, handlers.DownloadFileHandler)
		}

		// With middleware (error, auth or anonym)
		g = nfasthttp.NewRouterGroup(httppaths.GroupDownloader, r)
		g.Use(middlewareError, s.authMiddleware.AuthOrAnonym)
		{
			g.GET(httppaths.PathChannelAvatar, handlers.GetChannelAvatarHandler)
			g.HEAD(httppaths.PathChannelAvatar, handlers.GetChannelAvatarHandler)

			g.GET(httppaths.PathSettingsMenu, handlers.SettingsMenuHandler)

			g.GET(httppaths.PathStreamShortCode, handlers.StreamShortCodeHandler)
			g.HEAD(httppaths.PathStreamShortCode, handlers.StreamShortCodeHandler)
		}

	}

	// Short link
	{
		// With middleware (error, auth or anonym)
		g := nfasthttp.NewRouterGroup(s.shortLinkPrefix, r)
		g.Use(middlewareError, s.authMiddleware.AuthOrAnonym)
		{
			g.GET(httppaths.PathShortLink, handlers.ResolveShortLinkHandler)
			g.HEAD(httppaths.PathShortLink, handlers.ResolveShortLinkHandler)
		}
	}
}

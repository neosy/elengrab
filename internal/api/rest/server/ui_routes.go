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
		g.GET(httppaths.PathRegister, handlers.Downloader.AuthRegisterPageHandler)
		g.HEAD(httppaths.PathRegister, handlers.Downloader.AuthRegisterPageHandler)
		g.POST(httppaths.PathRegister, handlers.Downloader.AuthRegisterSubmitHandler)

		g.GET(httppaths.PathLogin, handlers.Downloader.AuthLoginPageHandler)
		g.HEAD(httppaths.PathLogin, handlers.Downloader.AuthLoginPageHandler)
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
			g.HEAD(httppaths.PathAccountMenu, handlers.Downloader.AccountMenuHandler)
		}

		// With middleware (error, auth optional)
		g = nfasthttp.NewRouterGroup(httppaths.GroupDownloader, r)
		g.Use(middlewareError, s.authMiddleware.AuthOptional)
		{
			g.GET(httppaths.PathHistory, handlers.Downloader.MediaHistoryHandler)
			g.HEAD(httppaths.PathHistory, handlers.Downloader.MediaHistoryHandler)
			g.GET(httppaths.PathEvents, handlers.Downloader.EventsStreamHandler)
			g.POST(httppaths.PathSearch, handlers.Downloader.SearchHandler)
		}

		// With middleware (error, auth or guest)
		g = nfasthttp.NewRouterGroup(httppaths.GroupDownloader, r)
		g.Use(middlewareError, s.authMiddleware.AuthOrGuest)
		{
			g.POST(httppaths.PathGrab, handlers.Downloader.ImportMediaByURLHandler)
			g.GET(httppaths.PathShareTarget, handlers.Downloader.ImportFromShareHandler)
			g.DELETE(httppaths.PathMediaItem, handlers.Downloader.MediaItemDeleteHandler)
			g.POST(httppaths.PathMediaItemDownloadRepeat, handlers.Downloader.RetryImportMediaHandler)
		}

		// With middleware (error, require auth mode)
		g = nfasthttp.NewRouterGroup(httppaths.GroupDownloader, r)
		g.Use(middlewareError, s.authMiddleware.RequireAuthMode)
		{
			g.GET(httppaths.PathMediaItemRow, handlers.Downloader.MediaItemRowHandler)
			g.HEAD(httppaths.PathMediaItemRow, handlers.Downloader.MediaItemRowHandler)
			g.GET(httppaths.PathMediaItemImage, handlers.Downloader.MediaItemImageHandler)
			g.HEAD(httppaths.PathMediaItemImage, handlers.Downloader.MediaItemImageHandler)
			g.GET(httppaths.PathMediaItemMenu, handlers.Downloader.MediaItemRowMenuHandler)
			g.POST(httppaths.PathMediaItemShortLink, handlers.Downloader.CreateMediaShareLinkHandler)
			g.GET(httppaths.PathMediaItemStream, handlers.Downloader.MediaItemStreamHandler)
			g.HEAD(httppaths.PathMediaItemStream, handlers.Downloader.MediaItemStreamHandler)
			g.GET(httppaths.PathMediaItemWatch, handlers.Downloader.WatchPageByDownloadIDHandler)
			g.HEAD(httppaths.PathMediaItemWatch, handlers.Downloader.WatchPageByDownloadIDHandler)

			g.GET(httppaths.PathDownloadFile, handlers.Downloader.DownloadFileHandler)
			g.HEAD(httppaths.PathDownloadFile, handlers.Downloader.DownloadFileHandler)
		}

		// With middleware (error, auth or anonym)
		g = nfasthttp.NewRouterGroup(httppaths.GroupDownloader, r)
		g.Use(middlewareError, s.authMiddleware.AuthOrAnonym)
		{
			g.GET(httppaths.PathChannelAvatar, handlers.Downloader.GetChannelAvatarHandler)
			g.HEAD(httppaths.PathChannelAvatar, handlers.Downloader.GetChannelAvatarHandler)

			g.GET(httppaths.PathStreamShortCode, handlers.Downloader.StreamShortCodeHandler)
			g.HEAD(httppaths.PathStreamShortCode, handlers.Downloader.StreamShortCodeHandler)
		}

	}

	// Short link
	{
		// With middleware (error, auth or anonym)
		g := nfasthttp.NewRouterGroup(s.shortLinkPrefix, r)
		g.Use(middlewareError, s.authMiddleware.AuthOrAnonym)
		{
			g.GET(httppaths.PathShortLink, handlers.Downloader.ResolveShortLinkHandler)
			g.HEAD(httppaths.PathShortLink, handlers.Downloader.ResolveShortLinkHandler)
		}
	}
}

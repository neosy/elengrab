package routes

import (
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
)

// registerUIDownloader register UI douwnloader routes.
func (r *routes) registerUIDownloader(handlers *downloader.DownloaderHandlers, shortLinkPrefix string) {
	middlewareError := r.middlewares.Error.ErrorHandler

	// Account with middleware (error, auth or anonym)
	g := nfasthttp.NewRouterGroup(httppaths.GroupAccount, r.router)
	g.Use(middlewareError, r.middlewares.Auth.AuthOrAnonym)
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
		g := nfasthttp.NewRouterGroup(httppaths.GroupDownloader, r.router)
		g.Use(middlewareError, r.middlewares.Auth.RequireAuth)
		{
			g.GET(httppaths.PathAccountMenu, handlers.AccountMenuHandler)
			g.HEAD(httppaths.PathAccountMenu, handlers.AccountMenuHandler)
		}

		// With middleware (error, auth optional)
		g = nfasthttp.NewRouterGroup(httppaths.GroupDownloader, r.router)
		g.Use(middlewareError, r.middlewares.Auth.AuthOptional)
		{
			g.HEAD(httppaths.PathHistory, handlers.MediaHistoryHandler)
			g.GET(httppaths.PathHistory, handlers.MediaHistoryHandler)

			g.GET(httppaths.PathEvents, handlers.EventsStreamHandler)

			g.POST(httppaths.PathSearch, handlers.SearchHandler)

			g.HEAD(httppaths.PathMediaItemRow, handlers.MediaItemRowHandler)
			g.GET(httppaths.PathMediaItemRow, handlers.MediaItemRowHandler)

			g.HEAD(httppaths.PathMediaItemImage, handlers.MediaItemImageHandler)
			g.GET(httppaths.PathMediaItemImage, handlers.MediaItemImageHandler)

			g.GET(httppaths.PathMediaItemMenu, handlers.MediaItemRowMenuHandler)

			g.HEAD(httppaths.PathMediaItemStream, handlers.MediaItemStreamHandler)
			g.GET(httppaths.PathMediaItemStream, handlers.MediaItemStreamHandler)

			g.HEAD(httppaths.PathMediaItemWatch, handlers.WatchPageByDownloadIDHandler)
			g.GET(httppaths.PathMediaItemWatch, handlers.WatchPageByDownloadIDHandler)

			g.HEAD(httppaths.PathDownloadFile, handlers.DownloadFileHandler)
			g.GET(httppaths.PathDownloadFile, handlers.DownloadFileHandler)

			g.HEAD(httppaths.PathMediaItemShortLink, handlers.GetMediaShareLinkHandler)
			g.GET(httppaths.PathMediaItemShortLink, handlers.GetMediaShareLinkHandler)
		}

		// With middleware (error, auth or guest)
		g = nfasthttp.NewRouterGroup(httppaths.GroupDownloader, r.router)
		g.Use(middlewareError, r.middlewares.Auth.AuthOrGuest)
		{
			g.POST(httppaths.PathGrab, handlers.ImportMediaByURLHandler)
			g.GET(httppaths.PathShareTarget, handlers.ImportFromShareHandler)
			g.DELETE(httppaths.PathMediaItem, handlers.MediaItemDeleteHandler)
			g.POST(httppaths.PathMediaItemDownloadRepeat, handlers.RetryImportMediaHandler)

			g.HEAD(httppaths.PathMediaItemEdit, handlers.EditMediaPageByDownloadIDHandler)
			g.GET(httppaths.PathMediaItemEdit, handlers.EditMediaPageByDownloadIDHandler)

			g.PATCH(httppaths.PathMediaItem, handlers.PatchMediaByDownloadIDHandler)

			g.POST(httppaths.PathMediaItemRefresh, handlers.RefreshMetadataByDownloadIDHandler)
		}

		// With middleware (error, require auth mode)
		g = nfasthttp.NewRouterGroup(httppaths.GroupDownloader, r.router)
		g.Use(middlewareError, r.middlewares.Auth.RequireAuthMode)
		{
			g.POST(httppaths.PathMediaItemShortLink, handlers.CreateMediaShareLinkHandler)
			g.DELETE(httppaths.PathMediaItemShortLink, handlers.DeleteMediaShareLinkHandler)
		}

		// With middleware (error, auth or anonym)
		g = nfasthttp.NewRouterGroup(httppaths.GroupDownloader, r.router)
		g.Use(middlewareError, r.middlewares.Auth.AuthOrAnonym)
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
		g := nfasthttp.NewRouterGroup(shortLinkPrefix, r.router)
		g.Use(middlewareError, r.middlewares.Auth.AuthOrAnonym)
		{
			g.GET(httppaths.PathShortLink, handlers.ResolveShortLinkHandler)
			g.HEAD(httppaths.PathShortLink, handlers.ResolveShortLinkHandler)
		}
	}
}

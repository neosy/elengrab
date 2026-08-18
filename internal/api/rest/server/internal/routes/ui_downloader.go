package routes

import (
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
)

// registerUIDownloader register UI douwnloader routes.
func (r *routes) registerUIDownloader(handlers *downloader.DownloaderHandlers, shortLinkPrefix string) {
	middlewareError := r.middlewares.Error.ErrorHandler

	// Downloader
	//group = r.Group(httppaths.GroupDownloader)
	{
		// With middleware (error, require auth)
		g := nfasthttp.NewRouterGroup(httppaths.DownloaderGroup, r.router)
		g.Use(middlewareError, r.middlewares.Auth.RequireAuth)
		{
			g.GET(httppaths.AccountMenuPath, handlers.AccountMenuHandler)
			g.HEAD(httppaths.AccountMenuPath, handlers.AccountMenuHandler)
		}

		// With middleware (error, auth optional)
		g = nfasthttp.NewRouterGroup(httppaths.DownloaderGroup, r.router)
		g.Use(middlewareError, r.middlewares.Auth.AuthOptional)
		{
			g.HEAD(httppaths.HistoryPath, handlers.MediaHistoryHandler)
			g.GET(httppaths.HistoryPath, handlers.MediaHistoryHandler)

			g.GET(httppaths.EventsPath, handlers.EventsStreamHandler)

			g.POST(httppaths.SearchPath, handlers.SearchHandler)

			g.HEAD(httppaths.MediaItemRowPath, handlers.MediaItemRowHandler)
			g.GET(httppaths.MediaItemRowPath, handlers.MediaItemRowHandler)

			g.HEAD(httppaths.MediaItemImagePath, handlers.MediaItemImageHandler)
			g.GET(httppaths.MediaItemImagePath, handlers.MediaItemImageHandler)

			g.GET(httppaths.MediaItemMenuPath, handlers.MediaItemRowMenuHandler)

			g.HEAD(httppaths.MediaItemStreamPath, handlers.MediaItemStreamHandler)
			g.GET(httppaths.MediaItemStreamPath, handlers.MediaItemStreamHandler)

			g.HEAD(httppaths.MediaItemWatchPath, handlers.WatchPageByDownloadIDHandler)
			g.GET(httppaths.MediaItemWatchPath, handlers.WatchPageByDownloadIDHandler)

			g.HEAD(httppaths.DownloadFilePath, handlers.DownloadFileHandler)
			g.GET(httppaths.DownloadFilePath, handlers.DownloadFileHandler)

			g.HEAD(httppaths.MediaItemShortLinkPath, handlers.GetShareLinkHandler)
			g.GET(httppaths.MediaItemShortLinkPath, handlers.GetShareLinkHandler)

			g.HEAD(httppaths.MediaItemWatchPositionPath, handlers.GetLastWatchPositionHandler)
			g.GET(httppaths.MediaItemWatchPositionPath, handlers.GetLastWatchPositionHandler)

			g.POST(httppaths.MediaItemReWatchTrackingPath, handlers.MediaItemWatchTrackingHandler)
		}

		// With middleware (error, auth or guest)
		g = nfasthttp.NewRouterGroup(httppaths.DownloaderGroup, r.router)
		g.Use(middlewareError, r.middlewares.Auth.AuthOrGuest)
		{
			g.POST(httppaths.GrabPath, handlers.ImportMediaByURLHandler)
			g.GET(httppaths.ShareTargetPath, handlers.ImportFromShareHandler)
			g.DELETE(httppaths.MediaItemPath, handlers.MediaItemDeleteHandler)
			g.POST(httppaths.MediaItemDownloadRepeatPath, handlers.RetryImportMediaHandler)

			g.HEAD(httppaths.MediaItemEditPath, handlers.EditMediaPageByDownloadIDHandler)
			g.GET(httppaths.MediaItemEditPath, handlers.EditMediaPageByDownloadIDHandler)

			g.PATCH(httppaths.MediaItemPath, handlers.PatchMediaByDownloadIDHandler)

			g.POST(httppaths.MediaItemRefreshPath, handlers.RefreshMetadataByDownloadIDHandler)
		}

		// With middleware (error, require auth mode)
		g = nfasthttp.NewRouterGroup(httppaths.DownloaderGroup, r.router)
		g.Use(middlewareError, r.middlewares.Auth.RequireAuthMode)
		{
			g.POST(httppaths.MediaItemShortLinkPath, handlers.CreateShareLinkHandler)
			g.DELETE(httppaths.MediaItemShortLinkPath, handlers.DeleteShareLinkHandler)
		}

		// With middleware (error, auth or anonym)
		g = nfasthttp.NewRouterGroup(httppaths.DownloaderGroup, r.router)
		g.Use(middlewareError, r.middlewares.Auth.AuthOrAnonym)
		{
			g.GET(httppaths.ChannelAvatarPath, handlers.GetChannelAvatarHandler)
			g.HEAD(httppaths.ChannelAvatarPath, handlers.GetChannelAvatarHandler)

			g.GET(httppaths.SettingsMenuPath, handlers.SettingsMenuHandler)

			g.GET(httppaths.StreamShortCodePath, handlers.StreamShortCodeHandler)
			g.HEAD(httppaths.StreamShortCodePath, handlers.StreamShortCodeHandler)
		}

	}

	// Short link
	{
		// With middleware (error, auth or anonym)
		g := nfasthttp.NewRouterGroup(shortLinkPrefix, r.router)
		g.Use(middlewareError, r.middlewares.Auth.AuthOrAnonym)
		{
			g.GET(httppaths.ShortLinkPath, handlers.ResolveShortLinkHandler)
			g.HEAD(httppaths.ShortLinkPath, handlers.ResolveShortLinkHandler)
		}
	}
}

package routes

import (
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
)

// registerUIDownloader register UI douwnloader routes.
func (r *routes) registerUIDownloader(handlers *downloader.DownloaderHandlers, shortLinkPrefix string) {
	// middlewareError renders an HTML error page only for failed GET requests
	// (other methods, including HEAD, keep the JSON error as is).
	middlewareError := r.middlewares.Error.ErrorHandler

	// Downloader
	//group = r.Group(httppaths.GroupDownloader)
	{
		// With middleware (error, require auth)
		// The error middleware renders an HTML error page only for failed GET requests;
		// errors of other methods (POST, PATCH, DELETE, HEAD) stay JSON.
		g := nfasthttp.NewRouterGroup(httppaths.DownloaderGroup, r.router)
		g.Use(middlewareError, r.middlewares.Auth.RequireAuth)
		{
			g.HEAD(httppaths.AccountMenuPath, handlers.AccountMenuHandler)
			g.GET(httppaths.AccountMenuPath, handlers.AccountMenuHandler)
		}

		// With middleware (auth optional)
		g = nfasthttp.NewRouterGroup(httppaths.DownloaderGroup, r.router)
		g.Use(r.middlewares.Auth.AuthOptional)
		{
			g.GET(httppaths.EventsPath, handlers.EventsStreamHandler)

			g.HEAD(httppaths.MediaItemStreamPath, handlers.MediaItemStreamHandler)
			g.GET(httppaths.MediaItemStreamPath, handlers.MediaItemStreamHandler)

			g.HEAD(httppaths.MediaItemShortLinkPath, handlers.GetShareLinkHandler)
			g.GET(httppaths.MediaItemShortLinkPath, handlers.GetShareLinkHandler)
			g.POST(httppaths.MediaItemShortLinkPath, handlers.CreateShareLinkHandler)

			g.HEAD(httppaths.MediaItemWatchPositionPath, handlers.GetLastWatchPositionHandler)
			g.GET(httppaths.MediaItemWatchPositionPath, handlers.GetLastWatchPositionHandler)
		}

		// With middleware (error, auth optional)
		// The error middleware renders an HTML error page only for failed GET requests;
		// errors of other methods (POST, PATCH, DELETE, HEAD) stay JSON.
		g = nfasthttp.NewRouterGroup(httppaths.DownloaderGroup, r.router)
		g.Use(middlewareError, r.middlewares.Auth.AuthOptional)
		{
			g.HEAD(httppaths.MediaItemsPath, handlers.MediaItemsHandler)
			g.GET(httppaths.MediaItemsPath, handlers.MediaItemsHandler)
			g.POST(httppaths.MediaItemsPath, handlers.MediaItemsHandler)

			g.POST(httppaths.SearchPath, handlers.SearchHandler)

			g.HEAD(httppaths.MediaItemRowPath, handlers.MediaItemRowHandler)
			g.GET(httppaths.MediaItemRowPath, handlers.MediaItemRowHandler)

			g.HEAD(httppaths.MediaItemImagePath, handlers.MediaItemImageHandler)
			g.GET(httppaths.MediaItemImagePath, handlers.MediaItemImageHandler)

			g.GET(httppaths.MediaItemMenuPath, handlers.MediaItemRowMenuHandler)

			g.HEAD(httppaths.MediaItemPath, handlers.MediaItemPageByDownloadIDHandler)
			g.GET(httppaths.MediaItemPath, handlers.MediaItemPageByDownloadIDHandler)

			g.HEAD(httppaths.DownloadFilePath, handlers.DownloadFileHandler)
			g.GET(httppaths.DownloadFilePath, handlers.DownloadFileHandler)

			g.POST(httppaths.MediaItemReWatchTrackingPath, handlers.MediaItemWatchTrackingHandler)
		}

		// With middleware (error, auth or guest)
		// The error middleware renders an HTML error page only for failed GET requests;
		// errors of other methods (POST, PATCH, DELETE, HEAD) stay JSON.
		g = nfasthttp.NewRouterGroup(httppaths.DownloaderGroup, r.router)
		g.Use(middlewareError, r.middlewares.Auth.AuthOrGuest)
		{
			g.POST(httppaths.GrabPath, handlers.QueueMediaDownloadByURLHandler)

			g.GET(httppaths.ShareTargetPath, handlers.ImportFromShareHandler)

			g.HEAD(httppaths.MediaItemEditPath, handlers.EditMediaPageByDownloadIDHandler)
			g.GET(httppaths.MediaItemEditPath, handlers.EditMediaPageByDownloadIDHandler)

			g.PATCH(httppaths.MediaItemPath, handlers.PatchMediaByDownloadIDHandler)
			g.DELETE(httppaths.MediaItemPath, handlers.MediaItemDeleteHandler)
			g.POST(httppaths.MediaItemDownloadRepeatPath, handlers.RetryImportMediaHandler)

			g.POST(httppaths.MediaItemRefreshPath, handlers.RefreshMetadataByDownloadIDHandler)

			g.HEAD(httppaths.UploadPath, handlers.UploadFileHandler)
			g.POST(httppaths.UploadPath, handlers.UploadFileHandler)

			g.HEAD(httppaths.UploadFilePath, handlers.QueueUploadedFileHandler)
			g.POST(httppaths.UploadFilePath, handlers.QueueUploadedFileHandler)
			g.DELETE(httppaths.UploadFilePath, handlers.DeleteUploadedFileHandler)
		}

		// With middleware (error, require auth mode)
		// The error middleware renders an HTML error page only for failed GET requests;
		// errors of other methods (POST, PATCH, DELETE, HEAD) stay JSON.
		g = nfasthttp.NewRouterGroup(httppaths.DownloaderGroup, r.router)
		g.Use(middlewareError, r.middlewares.Auth.RequireAuthMode)
		{
			g.DELETE(httppaths.MediaItemShortLinkPath, handlers.DeleteShareLinkHandler)
		}

		// With middleware (error, auth or anonym)
		// The error middleware renders an HTML error page only for failed GET requests;
		// errors of other methods (POST, PATCH, DELETE, HEAD) stay JSON.
		g = nfasthttp.NewRouterGroup(httppaths.DownloaderGroup, r.router)
		g.Use(middlewareError, r.middlewares.Auth.AuthOrAnonym)
		{
			g.GET(httppaths.SettingsMenuPath, handlers.SettingsMenuHandler)
		}

	}

	// Channels
	{
		// With middleware (error, auth or anonym)
		// The error middleware renders an HTML error page only for failed GET requests;
		// errors of other methods (POST, PATCH, DELETE, HEAD) stay JSON.
		g := nfasthttp.NewRouterGroup(httppaths.DownloaderGroup, r.router)
		g.Use(middlewareError, r.middlewares.Auth.AuthOrAnonym)
		{
			g.GET(httppaths.ChannelPath, handlers.ChannelPageHandler)
			g.HEAD(httppaths.ChannelPath, handlers.ChannelPageHandler)
		}
	}

	// User downloads
	{
		// With middleware (error, auth or anonym)
		g := nfasthttp.NewRouterGroup(httppaths.DownloaderGroup, r.router)
		g.Use(middlewareError, r.middlewares.Auth.AuthOrAnonym)
		{
			g.GET(httppaths.UserDownloadsPath, handlers.UserDownloadsPageHandler)
			g.HEAD(httppaths.UserDownloadsPath, handlers.UserDownloadsPageHandler)
		}
	}

	// Short link
	{
		// With middleware (auth or anonym)
		g := nfasthttp.NewRouterGroup(shortLinkPrefix, r.router)
		g.Use(r.middlewares.Auth.AuthOrAnonym)
		{
			g.HEAD(httppaths.ShortLinkStreamPath, handlers.ShortLinkStreamHandler)
			g.GET(httppaths.ShortLinkStreamPath, handlers.ShortLinkStreamHandler)
		}

		// With middleware (error, auth or anonym)
		// The error middleware renders an HTML error page only for failed GET requests;
		// errors of other methods (POST, PATCH, DELETE, HEAD) stay JSON.
		g = nfasthttp.NewRouterGroup(shortLinkPrefix, r.router)
		g.Use(middlewareError, r.middlewares.Auth.AuthOrAnonym)
		{
			g.HEAD(httppaths.ShortLinkPath, handlers.ResolveShortLinkHandler)
			g.GET(httppaths.ShortLinkPath, handlers.ResolveShortLinkHandler)

			g.HEAD(httppaths.ShortLinkImagePath, handlers.ShortLinkImageHandler)
			g.GET(httppaths.ShortLinkImagePath, handlers.ShortLinkImageHandler)
		}
	}
}

package usecases

import (
	"context"
	"log/slog"
	"time"

	"github.com/neosy/elengrab/internal/app/services"
	"github.com/neosy/elengrab/internal/app/usecases/admin"
	"github.com/neosy/elengrab/internal/app/usecases/auth"
	authweb "github.com/neosy/elengrab/internal/app/usecases/auth_web"
	"github.com/neosy/elengrab/internal/app/usecases/downloader"
	"github.com/neosy/elengrab/internal/app/usecases/link"
	linkweb "github.com/neosy/elengrab/internal/app/usecases/link_web"
	"github.com/neosy/elengrab/internal/app/usecases/maintenance"
	"github.com/neosy/elengrab/internal/app/usecases/thumbnail"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/workerpool"
	"github.com/neosy/elengrab/internal/ports/persistence"
	pstorage "github.com/neosy/elengrab/internal/ports/storage"
)

type Dependencies struct {
	Repositories DepRepositories
	Storages     DepStorages
	Services     *services.Services

	// dispetchers
	DownloadDispetcher   workerpool.JobDispatcher
	OperationDispatcher  workerpool.JobDispatcher
	WatchEventDispatcher workerpool.JobDispatcher

	// Options
	AppName  string
	AppMode  dtypes.AppMode
	DemoMode bool

	AuthSessionTTL             time.Duration
	AuthSessionRefreshInterval time.Duration

	BaseURL string

	// Short link
	BaseShortURL    string
	ShortCodeLength uint8
	ShortLinkTTL    time.Duration

	DatabaseBackupsDir  string
	DatabaseBackupsKeep int

	DeleteDuplicatesUniquenessScope dtypes.UniquenessScope

	LogoUpdateInterval    time.Duration
	ChannelUpdateInterval time.Duration

	DefaultAdminLogin    string
	DefaultAdminPassword string
}

type DepRepositories struct {
	Repositories persistence.Repositories

	User        persistence.UserRepository
	Role        persistence.RoleRepository
	UserRole    persistence.UserRoleRepository
	UserSession persistence.UserSessionRepository

	MediaDownload         persistence.MediaDownloadRepository
	DownloadTask          persistence.DownloadTaskRepository
	DownloadDataMigration persistence.DownloadDataMigrationRepository

	MediaWatchEvent        persistence.MediaWatchEventRepository
	MediaUserWatchChunk    persistence.MediaUserWatchChunkRepository
	MediaUserWatchStat     persistence.MediaUserWatchStatRepository
	MediaWatchStat         persistence.MediaWatchStatRepository
	MediaUserWatchPosition persistence.MediaUserWatchPositionRepository

	YoutubeChannel persistence.YoutubeChannelRepository
	SiteLogo       persistence.SiteLogoRepository
	Thumbnail      persistence.ThumbnailRepository

	Link      persistence.LinkRepository
	LinkClick persistence.LinkClickRepository

	// in memory
	MediaDownloadCache          persistence.MediaDownloadCacheRepository
	DownloadStateCache          persistence.DownloadStateCacheRepository
	MediaUserWatchStatCache     persistence.MediaUserWatchStatCacheRepository
	MediaWatchStatCache         persistence.MediaWatchStatCacheRepository
	MediaUserWatchPositionCache persistence.MediaUserWatchPositionCacheRepository
	YoutubeChannelCache         persistence.YoutubeChannelCacheRepository
	SiteLogoCache               persistence.SiteLogoCacheRepository
	ThumbnailCache              persistence.ThumbnailCacheRepository
	ThumbnailFileCache          persistence.ThumbnailFileCacheRepository
}

type DepStorages struct {
	Thumbnails pstorage.ThumbnailsStorage
	Downloads  pstorage.DownloadsStorage
}

type Usecases struct {
	Auth        *auth.Auth
	AuthWeb     *authweb.AuthWeb
	Admin       *admin.Admin
	Downloader  *downloader.Downloader
	Maintenance *maintenance.Maintenance
	Link        *link.Link
	LinkWeb     *linkweb.LinkWeb
	Thumbnail   *thumbnail.Thumbnail
}

func NewUsecases(ctx context.Context, logger *slog.Logger, deps *Dependencies) *Usecases {
	auth := auth.NewAuth(
		logger,
		deps.Repositories.User,
		deps.Repositories.Role,
		deps.Repositories.UserRole,
		deps.Repositories.UserSession,
		auth.WithSessionTTL(deps.AuthSessionTTL),
		auth.WithSessionRefreshInterval(deps.AuthSessionRefreshInterval),
	)

	link := link.NewLink(
		logger,
		deps.Repositories.Link,
		deps.Repositories.LinkClick,
		link.WithBaseURL(deps.BaseShortURL),
	)

	thumbnail := thumbnail.NewThumbnail(
		logger,
		deps.Repositories.Thumbnail,
		deps.Repositories.ThumbnailCache,
		deps.Repositories.ThumbnailFileCache,
		deps.Storages.Thumbnails,
	)

	return &Usecases{
		Auth: auth,
		AuthWeb: authweb.NewAuthWeb(
			logger,
			auth,
			deps.DefaultAdminLogin,
			deps.DefaultAdminPassword,
		),
		Admin: admin.NewAdmin(logger, auth),
		Downloader: downloader.NewDownloader(
			ctx,
			logger,

			// repositories
			deps.Repositories.MediaDownload,
			deps.Repositories.DownloadTask,
			deps.Repositories.DownloadDataMigration,
			deps.Repositories.MediaWatchEvent,
			deps.Repositories.MediaUserWatchChunk,
			deps.Repositories.MediaUserWatchStat,
			deps.Repositories.MediaWatchStat,
			deps.Repositories.MediaUserWatchPosition,
			deps.Repositories.YoutubeChannel,
			deps.Repositories.SiteLogo,

			// in memory
			deps.Repositories.MediaDownloadCache,
			deps.Repositories.DownloadStateCache,
			deps.Repositories.MediaUserWatchStatCache,
			deps.Repositories.MediaWatchStatCache,
			deps.Repositories.MediaUserWatchPositionCache,
			deps.Repositories.YoutubeChannelCache,
			deps.Repositories.SiteLogoCache,

			// storages
			deps.Storages.Downloads,

			// dispetchers
			deps.DownloadDispetcher,
			deps.OperationDispatcher,
			deps.WatchEventDispatcher,

			// usecases
			thumbnail,

			// services
			deps.Services.Downloader,
			deps.Services.FFMpeg,
			auth,

			// options
			deps.DemoMode,
			deps.AppMode,
			deps.DeleteDuplicatesUniquenessScope,
			deps.LogoUpdateInterval,
			deps.ChannelUpdateInterval,
		),
		Maintenance: maintenance.NewMaintenance(
			logger,
			deps.Repositories.Repositories,

			// options
			deps.AppName,
			deps.DatabaseBackupsDir,
			deps.DatabaseBackupsKeep,
		),
		Link: link,
		LinkWeb: linkweb.NewLinkWeb(
			logger,
			linkweb.Dependencies{Link: link},
			linkweb.Options{
				BaseShortURL:    deps.BaseShortURL,
				ShortCodeLength: deps.ShortCodeLength,
				LinkTTL:         deps.ShortLinkTTL,
			},
		),
		Thumbnail: thumbnail,
	}
}

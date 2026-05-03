package usecases

import (
	"context"
	"log/slog"
	"time"

	"github.com/neosy/elengrab/internal/app/services"
	"github.com/neosy/elengrab/internal/app/usecases/auth"
	authweb "github.com/neosy/elengrab/internal/app/usecases/auth_web"
	"github.com/neosy/elengrab/internal/app/usecases/downloader"
	"github.com/neosy/elengrab/internal/app/usecases/link"
	linkweb "github.com/neosy/elengrab/internal/app/usecases/link_web"
	"github.com/neosy/elengrab/internal/app/usecases/maintenance"
	"github.com/neosy/elengrab/internal/app/usecases/thumbnail"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	nworkerpool "github.com/neosy/elengrab/internal/pkg/workerpool"
	"github.com/neosy/elengrab/internal/ports/persistence"
	pstorage "github.com/neosy/elengrab/internal/ports/storage"
)

type Dependencies struct {
	Repositories DepRepositories
	Storages     DepStorages
	Services     *services.Services

	// dispetchers
	DownloadDispetcher nworkerpool.JobDispatcher

	// Options
	AppName  string
	AppMode  dtypes.AppMode
	DemoMode bool

	AuthSessionTTL             time.Duration
	AuthSessionRefreshInterval time.Duration

	BaseURL string

	BaseShortURL    string
	ShortCodeLength uint8

	DatabaseBackupsDir  string
	DatabaseBackupsKeep int

	DeleteDuplicatesScope dtypes.UniquenessScope

	LogoUpdateInterval    time.Duration
	ChannelUpdateInterval time.Duration

	DefaultAdminLogin    string
	DefaultAdminPassword string
}

type DepRepositories struct {
	Database persistence.Database

	User        persistence.UserRepository
	Role        persistence.RoleRepository
	UserRole    persistence.UserRoleRepository
	UserSession persistence.UserSessionRepository

	File                  persistence.FileRepository
	DownloadTask          persistence.DownloadTaskRepository
	DownloadDataMigration persistence.DownloadDataMigrationRepository

	YoutubeChannel persistence.YoutubeChannelRepository
	SiteLogo       persistence.SiteLogoRepository
	Thumbnail      persistence.ThumbnailRepository

	Link      persistence.LinkRepository
	LinkClick persistence.LinkClickRepository

	// in memory
	DownloadStateCache  persistence.DownloadStateCacheRepository
	YoutubeChannelCache persistence.YoutubeChannelCacheRepository
	SiteLogoCache       persistence.SiteLogoCacheRepository
}

type DepStorages struct {
	Thumbnails pstorage.ThumbnailsStorage
	Downloads  pstorage.DownloadsStorage
}

type Usecases struct {
	Auth        *auth.Auth
	AuthWeb     *authweb.AuthWeb
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
		link.LinkOptionBaseURL(deps.BaseShortURL),
	)

	thumbnail := thumbnail.NewThumbnail(
		logger,
		deps.Repositories.Thumbnail,
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
		Downloader: downloader.NewDownloader(
			ctx,
			logger,

			// repositories
			deps.Repositories.File,
			deps.Repositories.DownloadTask,
			deps.Repositories.DownloadDataMigration,
			deps.Repositories.YoutubeChannel,
			deps.Repositories.SiteLogo,
			deps.Repositories.Thumbnail,

			// in memory
			deps.Repositories.DownloadStateCache,
			deps.Repositories.YoutubeChannelCache,
			deps.Repositories.SiteLogoCache,

			// storages
			deps.Storages.Downloads,

			// dispetchers
			deps.DownloadDispetcher,

			// usecases
			thumbnail,

			// services
			deps.Services.YouTubeDownloader,

			// options
			deps.DemoMode,
			deps.AppMode,
			deps.DeleteDuplicatesScope,
			deps.LogoUpdateInterval,
			deps.ChannelUpdateInterval,
		),
		Maintenance: maintenance.NewMaintenance(
			logger,
			deps.Repositories.Database,

			// options
			deps.AppName,
			deps.DatabaseBackupsDir,
			deps.DatabaseBackupsKeep,
		),
		Link:      link,
		LinkWeb:   linkweb.NewLinkWeb(logger, link, deps.BaseShortURL, deps.ShortCodeLength),
		Thumbnail: thumbnail,
	}
}

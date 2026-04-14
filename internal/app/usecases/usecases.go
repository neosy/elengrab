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
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	nworkerpool "github.com/neosy/elengrab/internal/pkg/workerpool"
	"github.com/neosy/elengrab/internal/ports/persistence"
)

type Dependencies struct {
	Repositories DepRepositories
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

	DownloadsDir string

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

	File           persistence.FileRepository
	DownloadTask   persistence.DownloadTaskRepository
	YoutubeChannel persistence.YoutubeChannelRepository
	SiteLogo       persistence.SiteLogoRepository

	Link      persistence.LinkRepository
	LinkClick persistence.LinkClickRepository

	// in memory
	DownloadStateCache  persistence.DownloadStateCacheRepository
	YoutubeChannelCache persistence.YoutubeChannelCacheRepository
	SiteLogoCache       persistence.SiteLogoCacheRepository
}

type Usecases struct {
	Auth        *auth.Auth
	AuthWeb     *authweb.AuthWeb
	Downloader  *downloader.Downloader
	Maintenance *maintenance.Maintenance
	Link        *link.Link
	LinkWeb     *linkweb.LinkWeb
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
			deps.Repositories.YoutubeChannel,
			deps.Repositories.SiteLogo,

			// in memory
			deps.Repositories.DownloadStateCache,
			deps.Repositories.YoutubeChannelCache,
			deps.Repositories.SiteLogoCache,

			// dispetchers
			deps.DownloadDispetcher,

			// services
			deps.Services.YouTubeDownloader,

			// options
			deps.DemoMode,
			deps.DownloadsDir,
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
		Link:    link,
		LinkWeb: linkweb.NewLinkWeb(logger, link, deps.BaseShortURL, deps.ShortCodeLength),
	}
}

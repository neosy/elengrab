package usecases

import (
	"context"
	"log/slog"
	"time"

	"github.com/neosy/elengrab/internal/app/services"
	"github.com/neosy/elengrab/internal/app/usecases/auth"
	"github.com/neosy/elengrab/internal/app/usecases/authweb"
	"github.com/neosy/elengrab/internal/app/usecases/downloader"
	"github.com/neosy/elengrab/internal/app/usecases/maintenance"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/nworkerpool"
	"github.com/neosy/elengrab/internal/ports/persistence"
)

type Dependencies struct {
	Repositories DepRepositories
	Services     *services.Services

	// dispetchers
	DownloadDispetcher nworkerpool.JobDispatcher

	// Options
	AppName      string
	DemoMode     bool
	DownloadsDir string

	DatabaseBackupsDir  string
	DatabaseBackupsKeep int

	AppMode               dtypes.AppMode
	DeleteDuplicatesScope dtypes.UniquenessScope

	LogoUpdateInterval    time.Duration
	ChannelUpdateInterval time.Duration

	DefaultAdminLogin    string
	DefaultAdminPassword string
}

type DepRepositories struct {
	Database persistence.Database

	File           persistence.FileRepository
	DownloadTask   persistence.DownloadTaskRepository
	YoutubeChannel persistence.YoutubeChannelRepository
	SiteLogo       persistence.SiteLogoRepository

	User        persistence.UserRepository
	Role        persistence.RoleRepository
	UserRole    persistence.UserRoleRepository
	UserSession persistence.UserSessionRepository

	// in memory
	DownloadStateCache  persistence.DownloadStateCacheRepository
	YoutubeChannelCache persistence.YoutubeChannelCacheRepository
	SiteLogoCache       persistence.SiteLogoCacheRepository
}

type Usecases struct {
	Downloader  *downloader.Downloader
	Maintenance *maintenance.Maintenance
	Auth        *auth.Auth
	AuthWeb     *authweb.AuthWeb
}

func NewUsecases(ctx context.Context, logger *slog.Logger, deps *Dependencies) *Usecases {
	auth := auth.NewAuth(
		logger,
		deps.Repositories.User,
		deps.Repositories.Role,
		deps.Repositories.UserRole,
		deps.Repositories.UserSession,
	)
	return &Usecases{
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
		Auth: auth,
		AuthWeb: authweb.NewAuthWeb(
			logger,
			auth,
			deps.DefaultAdminLogin,
			deps.DefaultAdminPassword,
		),
	}
}

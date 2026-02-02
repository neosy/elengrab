package usecases

import (
	"context"
	"log/slog"

	"github.com/neosy/elengrab/internal/app/services"
	"github.com/neosy/elengrab/internal/app/usecases/auth"
	ytdownloader "github.com/neosy/elengrab/internal/app/usecases/downloader"
	"github.com/neosy/elengrab/internal/app/usecases/maintenance"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/ports/persistence"
	"github.com/neosy/elengrab/pkg/nworkerpool"
)

type Dependencies struct {
	Repositories DepRepositories
	Services     *services.Services

	// dispetchers
	DownloadDispetcher nworkerpool.JobDispatcher

	// Options
	AppName      string
	DownloadsDir string

	DatabaseBackupsDir  string
	DatabaseBackupsKeep int

	HistoryMode           dtypes.HistoryMode
	DeleteDuplicatesScope dtypes.UniquenessScope
}

type DepRepositories struct {
	Database persistence.Database

	File           persistence.FileRepository
	DownloadTask   persistence.DownloadTaskRepository
	YoutubeChannel persistence.YoutubeChannelRepository
	User           persistence.UserRepository
	UserSession    persistence.UserSessionRepository

	// in memory
	DownloadStateCache  persistence.DownloadStateRepository
	YoutubeChannelCache persistence.YoutubeChannelRepository
}

type Usecases struct {
	Downloader  *ytdownloader.YouTubeDownloader
	Maintenance *maintenance.Maintenance
	Auth        *auth.Auth
}

func NewUsecases(ctx context.Context, logger *slog.Logger, deps *Dependencies) *Usecases {
	return &Usecases{
		Downloader: ytdownloader.NewYouTubeDownloader(
			ctx,
			logger,

			// repositories
			deps.Repositories.File,
			deps.Repositories.DownloadTask,
			deps.Repositories.YoutubeChannel,

			// in memory
			deps.Repositories.DownloadStateCache,
			deps.Repositories.YoutubeChannelCache,

			// dispetchers
			deps.DownloadDispetcher,

			// services
			deps.Services.YouTubeDownloader,

			// options
			deps.DownloadsDir,
			deps.HistoryMode,
			deps.DeleteDuplicatesScope,
		),
		Maintenance: maintenance.NewMaintenance(
			logger,
			deps.Repositories.Database,

			// options
			deps.AppName,
			deps.DatabaseBackupsDir,
			deps.DatabaseBackupsKeep,
		),
		Auth: auth.NewAuth(
			logger,
			deps.Repositories.User,
			deps.Repositories.UserSession,
		),
	}
}

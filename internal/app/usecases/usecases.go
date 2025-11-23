package usecases

import (
	"log/slog"

	ytdownloader "github.com/neosy/elengrab/internal/app/usecases/youtube_downloader"
	"github.com/neosy/elengrab/internal/ports/persistence"
	"github.com/neosy/elengrab/internal/services"
	"github.com/neosy/elengrab/pkg/nworkerpool"
)

type Dependencies struct {
	Repositories DepRepositories
	Services     *services.Services

	// dispetchers
	DownloadDispetcher nworkerpool.JobDispatcher

	// Options
	DownloadsDir string
	LoadHistory  bool
}

type DepRepositories struct {
	File          persistence.FileRepository
	DownloadTask  persistence.DownloadTaskRepository
	DownloadState persistence.DownloadStateRepository
}

type Usecases struct {
	Downloader *ytdownloader.YouTubeDownloader
}

func NewUsecases(logger *slog.Logger, deps *Dependencies) *Usecases {
	return &Usecases{
		Downloader: ytdownloader.NewYouTubeDownloader(
			logger,

			// repositories
			deps.Repositories.File,
			deps.Repositories.DownloadTask,
			deps.Repositories.DownloadState,

			// dispetchers
			deps.DownloadDispetcher,

			// services
			deps.Services.YouTubeDownloader,

			// options
			deps.DownloadsDir,
			deps.LoadHistory,
		),
	}
}

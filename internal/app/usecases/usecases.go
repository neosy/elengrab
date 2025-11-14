package usecases

import (
	"log/slog"

	ytdownloader "github.com/neosy/elengrab/internal/app/usecases/youtube_downloader"
	"github.com/neosy/elengrab/internal/ports/persistence"
	"github.com/neosy/elengrab/internal/services"
	"github.com/neosy/elengrab/pkg/workerpool"
)

type Dependencies struct {
	Repositories DepRepositories
	Services     *services.Services

	// dispetchers
	DownloadDispetcher workerpool.JobDispatcher

	// Options
	DownloadsDir string
}

type DepRepositories struct {
	File         persistence.FileRepository
	DownloadTask persistence.DownloadTaskRepository
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

			// dispetchers
			deps.DownloadDispetcher,

			// services
			deps.Services.YouTubeDownloader,

			// options
			deps.DownloadsDir,
		),
	}
}

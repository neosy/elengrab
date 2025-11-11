package usecases

import (
	"log/slog"

	ucdownloader "github.com/neosy/elengrab/internal/app/usecases/youtube_downloader"
	"github.com/neosy/elengrab/internal/ports/persistence"
	"github.com/neosy/elengrab/internal/services"
)

type Dependencies struct {
	Repositories DepRepositories
	Services     *services.Services

	// Options
	DownloadsDir string
}

type DepRepositories struct {
	File         persistence.FileRepository
	DownloadTask persistence.DownloadTaskRepository
}

type Usecases struct {
	Downloader *ucdownloader.YouTubeDownloader
}

func NewUsecases(logger *slog.Logger, deps *Dependencies) *Usecases {
	return &Usecases{
		Downloader: ucdownloader.NewYouTubeDownloader(
			logger,

			// repositories
			deps.Repositories.File,
			deps.Repositories.DownloadTask,

			// services
			deps.Services.YouTubeDownloader,

			// options
			deps.DownloadsDir,
		),
	}
}

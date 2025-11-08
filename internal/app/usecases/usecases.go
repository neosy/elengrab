package usecases

import (
	"log/slog"

	ucdownloader "github.com/neosy/elengrab/internal/app/usecases/youtube_downloader"
	"github.com/neosy/elengrab/internal/services"
)

type Dependencies struct {
	Services *services.Services

	// Options
	DownloadsDir string
}

type Usecases struct {
	Downloader *ucdownloader.YouTubeDownloader
}

func NewUsecases(logger *slog.Logger, deps *Dependencies) *Usecases {
	return &Usecases{
		Downloader: ucdownloader.NewYouTubeDownloader(logger, deps.Services.YouTubeDownloader, deps.DownloadsDir),
	}
}

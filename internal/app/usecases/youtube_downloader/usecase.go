package ucdownloader

import (
	"log/slog"

	pservices "github.com/neosy/elengrab/internal/ports/services"
)

type YouTubeDownloader struct {
	logger *slog.Logger

	// services
	downloaderSrv pservices.YouTubeDownloader

	// Options
	downloadsDir string
}

func NewYouTubeDownloader(logger *slog.Logger, downloaderSrv pservices.YouTubeDownloader, downloadsDir string) *YouTubeDownloader {
	return &YouTubeDownloader{
		logger: logger,

		// services
		downloaderSrv: downloaderSrv,

		// Options
		downloadsDir: downloadsDir,
	}
}

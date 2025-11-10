package ucdownloader

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/ports/persistence"
	pservices "github.com/neosy/elengrab/internal/ports/services"
)

type YouTubeDownloader struct {
	logger *slog.Logger

	// repositories
	fileRep persistence.FileRepository

	// services
	downloaderSrv pservices.YouTubeDownloader

	// Options
	downloadsDir string
}

func NewYouTubeDownloader(
	logger *slog.Logger,

	// repositories
	fileRep persistence.FileRepository,

	// services
	downloaderSrv pservices.YouTubeDownloader,

	// options
	downloadsDir string,
) *YouTubeDownloader {
	return &YouTubeDownloader{
		logger: logger,

		// repositories
		fileRep: fileRep,

		// services
		downloaderSrv: downloaderSrv,

		// Options
		downloadsDir: downloadsDir,
	}
}

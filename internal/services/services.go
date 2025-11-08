package services

import (
	"log/slog"

	pservices "github.com/neosy/elengrab/internal/ports/services"
	ytdlpsrv "github.com/neosy/elengrab/internal/services/ytdlp"
)

type Dependencies struct {
	BinDir       string
	DownloadsDir string
}

type Services struct {
	YouTubeDownloader pservices.YouTubeDownloader
}

func New(logger *slog.Logger, deps *Dependencies) (*Services, error) {
	downloader, err := ytdlpsrv.NewYtDlpService(logger, deps.BinDir, deps.DownloadsDir)
	if err != nil {
		return nil, err
	}

	return &Services{
		YouTubeDownloader: downloader,
	}, nil
}

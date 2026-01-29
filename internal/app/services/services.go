package services

import (
	"log/slog"

	ytdlpsrv "github.com/neosy/elengrab/internal/app/services/ytdlp"
	pservices "github.com/neosy/elengrab/internal/ports/services"
)

type Dependencies struct {
	DownloaderBinDir string
	DownloadsDir     string

	YtDlpOptions *ytdlpsrv.Options
}

type Services struct {
	YouTubeDownloader pservices.Downloader
}

func New(logger *slog.Logger, deps *Dependencies) (*Services, error) {
	downloader, err := ytdlpsrv.NewYtDlpService(logger, deps.DownloaderBinDir, deps.DownloadsDir, deps.YtDlpOptions)
	if err != nil {
		return nil, err
	}

	return &Services{
		YouTubeDownloader: downloader,
	}, nil
}

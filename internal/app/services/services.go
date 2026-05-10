package services

import (
	"log/slog"

	ffmpegsrv "github.com/neosy/elengrab/internal/app/services/ffmpeg"
	ytdlpsrv "github.com/neosy/elengrab/internal/app/services/ytdlp"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/dto"
	pservices "github.com/neosy/elengrab/internal/ports/services"
	pstorage "github.com/neosy/elengrab/internal/ports/storage"
)

type Dependencies struct {
	DownloaderBinDir string
	Storage          pstorage.DownloadsStorage
	YtDlpOptions     []dto.Option
}

type Services struct {
	Downloader pservices.Downloader
	FFMpeg     pservices.FFMpeg
}

func New(logger *slog.Logger, deps *Dependencies) (*Services, error) {
	ffmpeg, err := ffmpegsrv.NewFFmpegService(logger, "")
	if err != nil {
		return nil, err
	}

	downloader, err := ytdlpsrv.NewYtDlpService(
		logger,
		deps.DownloaderBinDir,
		deps.Storage,
		ffmpeg,
		deps.YtDlpOptions...,
	)
	if err != nil {
		return nil, err
	}

	return &Services{
		Downloader: downloader,
		FFMpeg:     ffmpeg,
	}, nil
}

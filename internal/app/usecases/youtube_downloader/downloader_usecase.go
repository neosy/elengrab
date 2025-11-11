package ucdownloader

import (
	"log/slog"

	downloadtask "github.com/neosy/elengrab/internal/app/usecases/youtube_downloader/internal/download_task"
	fileuc "github.com/neosy/elengrab/internal/app/usecases/youtube_downloader/internal/file"
	filestatus "github.com/neosy/elengrab/internal/app/usecases/youtube_downloader/internal/file_status"
	"github.com/neosy/elengrab/internal/ports/persistence"
	pservices "github.com/neosy/elengrab/internal/ports/services"
)

type YouTubeDownloader struct {
	logger *slog.Logger

	// repositories
	// fileRep         persistence.FileRepository
	// downloadTaskRep persistence.DownloadTaskRepository

	// internal
	file         *fileuc.File
	fileStatus   *filestatus.FileStatus
	downloadTask *downloadtask.DownloadTask

	// services
	downloaderSrv pservices.YouTubeDownloader

	// Options
	downloadsDir string
}

func NewYouTubeDownloader(
	logger *slog.Logger,

	// repositories
	fileRep persistence.FileRepository,
	downloadTaskRep persistence.DownloadTaskRepository,

	// services
	downloaderSrv pservices.YouTubeDownloader,

	// options
	downloadsDir string,
) *YouTubeDownloader {
	downloadTask := downloadtask.NewDownloadTask(logger, downloadTaskRep)

	return &YouTubeDownloader{
		logger: logger,

		// repositories
		// fileRep: fileRep,

		// internal
		file:         fileuc.NewFile(logger, fileRep),
		fileStatus:   filestatus.NewOrderStatus(logger, fileRep, downloadTask),
		downloadTask: downloadTask,

		// services
		downloaderSrv: downloaderSrv,

		// Options
		downloadsDir: downloadsDir,
	}
}

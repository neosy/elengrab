package ytdownloader

import (
	"log/slog"

	dltask "github.com/neosy/elengrab/internal/app/usecases/youtube_downloader/internal/download_task"
	dltasktatus "github.com/neosy/elengrab/internal/app/usecases/youtube_downloader/internal/download_task_status"
	fileuc "github.com/neosy/elengrab/internal/app/usecases/youtube_downloader/internal/file"
	filestatus "github.com/neosy/elengrab/internal/app/usecases/youtube_downloader/internal/file_status"
	"github.com/neosy/elengrab/internal/ports/persistence"
	pservices "github.com/neosy/elengrab/internal/ports/services"
	"github.com/neosy/elengrab/pkg/workerpool"
)

type YouTubeDownloader struct {
	logger *slog.Logger

	// repositories
	// fileRep         persistence.FileRepository
	// downloadTaskRep persistence.DownloadTaskRepository

	// dispetchers
	dlDispetcher workerpool.JobDispatcher

	// internal
	file         *fileuc.File
	dlTask       *dltask.DownloadTask
	fileStatus   *filestatus.FileStatus
	dlTaskStatus *dltasktatus.DownloadTaskStatus

	// services
	downloaderSrv pservices.YouTubeDownloader

	// Options
	downloadsDir string
}

func NewYouTubeDownloader(
	logger *slog.Logger,

	// repositories
	fileRep persistence.FileRepository,
	dlTaskRep persistence.DownloadTaskRepository,

	// dispetchers
	dlDispetcher workerpool.JobDispatcher,

	// services
	downloaderSrv pservices.YouTubeDownloader,

	// options
	downloadsDir string,
) *YouTubeDownloader {
	dlTask := dltask.NewDownloadTask(logger, dlTaskRep)
	file := fileuc.NewFile(logger, fileRep, dlTask)
	dlTaskStatus := dltasktatus.NewDownloadTaskStatus(logger, dlTaskRep, dlTask)

	return &YouTubeDownloader{
		logger: logger,

		// repositories
		// fileRep: fileRep,

		// dispetchers
		dlDispetcher: dlDispetcher,

		// internal
		file:         file,
		dlTask:       dlTask,
		fileStatus:   filestatus.NewFileStatus(logger, fileRep, file, dlTask, dlTaskStatus),
		dlTaskStatus: dlTaskStatus,

		// services
		downloaderSrv: downloaderSrv,

		// Options
		downloadsDir: downloadsDir,
	}
}

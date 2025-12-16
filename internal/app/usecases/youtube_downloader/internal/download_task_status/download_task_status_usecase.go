package dltaskstatus

import (
	"log/slog"

	downloadtask "github.com/neosy/elengrab/internal/app/usecases/youtube_downloader/internal/download_task"
	statussetter "github.com/neosy/elengrab/internal/app/usecases/youtube_downloader/internal/download_task_status/status_setter"
)

type DownloadTaskStatus struct {
	logger *slog.Logger

	// internal
	statusSetter *statussetter.DownloadTaskStatusSetter

	// usecases
	downloadTask *downloadtask.DownloadTask
}

func NewDownloadTaskStatus(
	logger *slog.Logger,

	// usecases
	downloadTask *downloadtask.DownloadTask,
) *DownloadTaskStatus {
	return &DownloadTaskStatus{
		logger: logger,

		// internal
		statusSetter: statussetter.NewFDownloadTaskStatusSetter(logger),

		// usecases
		downloadTask: downloadTask,
	}
}

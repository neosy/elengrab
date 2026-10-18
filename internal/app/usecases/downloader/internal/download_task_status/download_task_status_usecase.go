package dltaskstatus

import (
	"log/slog"

	downloadtask "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/download_task"
	statussetter "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/download_task_status/status_setter"
)

type DownloadTaskStatus struct {
	logger *slog.Logger

	// internal
	statusSetter *statussetter.StatusSetter

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
		statusSetter: statussetter.NewStatusSetter(logger),

		// usecases
		downloadTask: downloadTask,
	}
}

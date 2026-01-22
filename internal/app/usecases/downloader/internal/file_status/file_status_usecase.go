package filestatus

import (
	"log/slog"

	dltask "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/download_task"
	dltasktatus "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/download_task_status"
	fileuc "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/file"
	statussetter "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/file_status/status_setter"
)

type FileStatus struct {
	logger *slog.Logger

	// internal
	statusSetter *statussetter.FileStatusSetter

	// usecases
	file         *fileuc.File
	dlTask       *dltask.DownloadTask
	dlTaskStatus *dltasktatus.DownloadTaskStatus
}

func NewFileStatus(
	logger *slog.Logger,

	// usecases
	file *fileuc.File,
	dlTask *dltask.DownloadTask,
	dlTaskStatus *dltasktatus.DownloadTaskStatus,
) *FileStatus {
	return &FileStatus{
		logger: logger,

		// internal
		statusSetter: statussetter.NewFileStatusSetter(logger),

		// usecases
		file:         file,
		dlTask:       dlTask,
		dlTaskStatus: dlTaskStatus,
	}
}

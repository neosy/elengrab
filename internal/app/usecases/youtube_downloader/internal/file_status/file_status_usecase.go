package filestatus

import (
	"log/slog"

	dltask "github.com/neosy/elengrab/internal/app/usecases/youtube_downloader/internal/download_task"
	dltasktatus "github.com/neosy/elengrab/internal/app/usecases/youtube_downloader/internal/download_task_status"
	fileuc "github.com/neosy/elengrab/internal/app/usecases/youtube_downloader/internal/file"
	statussetter "github.com/neosy/elengrab/internal/app/usecases/youtube_downloader/internal/file_status/status_setter"
	"github.com/neosy/elengrab/internal/ports/persistence"
)

type FileStatus struct {
	logger  *slog.Logger
	fileRep persistence.FileRepository

	// internal
	statusSetter *statussetter.FileStatusSetter

	// usecases
	file         *fileuc.File
	dlTask       *dltask.DownloadTask
	dlTaskStatus *dltasktatus.DownloadTaskStatus
}

func NewFileStatus(
	logger *slog.Logger,

	// repositories
	fileRep persistence.FileRepository,

	// usecases
	file *fileuc.File,
	dlTask *dltask.DownloadTask,
	dlTaskStatus *dltasktatus.DownloadTaskStatus,
) *FileStatus {
	return &FileStatus{
		logger:  logger,
		fileRep: fileRep,

		// internal
		statusSetter: statussetter.NewFileStatusSetter(logger),

		// usecases
		file:         file,
		dlTask:       dlTask,
		dlTaskStatus: dlTaskStatus,
	}
}

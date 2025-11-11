package filestatus

import (
	"log/slog"

	downloadtask "github.com/neosy/elengrab/internal/app/usecases/youtube_downloader/internal/download_task"
	statussetter "github.com/neosy/elengrab/internal/app/usecases/youtube_downloader/internal/file_status/status_setter"
	"github.com/neosy/elengrab/internal/ports/persistence"
)

type FileStatus struct {
	logger  *slog.Logger
	fileRep persistence.FileRepository

	// internal
	statusSetter *statussetter.FileStatusSetter

	// usecases
	downloadTask *downloadtask.DownloadTask
}

func NewOrderStatus(
	logger *slog.Logger,

	// repositories
	fileRep persistence.FileRepository,

	// usecases
	downloadTask *downloadtask.DownloadTask,
) *FileStatus {
	return &FileStatus{
		logger:  logger,
		fileRep: fileRep,

		// internal
		statusSetter: statussetter.NewFileStatusSetter(logger),

		// usecases
		downloadTask: downloadTask,
	}
}

package fileuc

import (
	"log/slog"

	dltask "github.com/neosy/elengrab/internal/app/usecases/youtube_downloader/internal/download_task"
	"github.com/neosy/elengrab/internal/ports/persistence"
)

type File struct {
	logger *slog.Logger

	// repositories
	fileRep persistence.FileRepository

	// usecases
	dlTask *dltask.DownloadTask
}

func NewFile(
	logger *slog.Logger,

	// repositories
	fileRep persistence.FileRepository,

	// usecases
	dlTask *dltask.DownloadTask,
) *File {
	return &File{
		logger:  logger,
		fileRep: fileRep,
		dlTask:  dlTask,
	}
}

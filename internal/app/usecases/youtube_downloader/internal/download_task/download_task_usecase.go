package downloadtask

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

type DownloadTask struct {
	logger *slog.Logger

	// repositories
	taskRep persistence.DownloadTaskRepository
}

func NewDownloadTask(
	logger *slog.Logger,
	taskRep persistence.DownloadTaskRepository,
) *DownloadTask {
	return &DownloadTask{
		logger:  logger,
		taskRep: taskRep,
	}
}

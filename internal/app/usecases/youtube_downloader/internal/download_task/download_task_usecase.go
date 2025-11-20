package dltask

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

type DownloadTask struct {
	logger *slog.Logger

	// repositories
	TaskRep persistence.DownloadTaskRepository
}

func NewDownloadTask(
	logger *slog.Logger,
	taskRep persistence.DownloadTaskRepository,
) *DownloadTask {
	return &DownloadTask{
		logger:  logger,
		TaskRep: taskRep,
	}
}

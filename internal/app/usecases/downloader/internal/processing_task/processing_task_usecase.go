package processingtask

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

type ProcessingTask struct {
	logger *slog.Logger

	// repositories
	TaskRepo persistence.MediaProcessingTaskRepositoryFactory
}

func NewProcessingTask(
	logger *slog.Logger,
	taskRepo persistence.MediaProcessingTaskRepositoryFactory,
) *ProcessingTask {
	return &ProcessingTask{
		logger:   logger,
		TaskRepo: taskRepo,
	}
}

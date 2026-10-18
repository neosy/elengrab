package proctaskstatus

import (
	"log/slog"

	processingtask "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/processing_task"
	statussetter "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/processing_task_status/status_setter"
)

type ProcessingTaskStatus struct {
	logger *slog.Logger

	// internal
	statusSetter *statussetter.StatusSetter

	// usecases
	processingTask *processingtask.ProcessingTask
}

func NewProcessingTaskStatus(
	logger *slog.Logger,

	// usecases
	processingTask *processingtask.ProcessingTask,
) *ProcessingTaskStatus {
	return &ProcessingTaskStatus{
		logger: logger,

		// internal
		statusSetter: statussetter.NewStatusSetter(logger),

		// usecases
		processingTask: processingTask,
	}
}

package wjobs

import (
	"context"
	"log/slog"

	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

type deleteMissingFilesJob struct {
	logger *slog.Logger
	runner pworkers.MaintenanceRunner

	// options
	enableMoveUnmatchedFiles bool
}

func NewDeleteMissingFilesJob(logger *slog.Logger, runner pworkers.MaintenanceRunner, enableMoveUnmatchedFiles bool) *deleteMissingFilesJob {
	return &deleteMissingFilesJob{
		logger: logger,
		runner: runner,

		// options
		enableMoveUnmatchedFiles: enableMoveUnmatchedFiles,
	}
}

func (j *deleteMissingFilesJob) Execute(ctx context.Context) error {
	err := j.runner.DeleteMissingFiles(ctx, j.enableMoveUnmatchedFiles)
	j.logger.Debug("Job done", "name", "DeleteMissingFiles")
	return err
}

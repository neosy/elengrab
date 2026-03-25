package wjobs

import (
	"context"
	"log/slog"
	"time"

	uformat "github.com/neosy/elengrab/internal/pkg/utils/format"
	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

type deleteMissingFilesJob struct {
	logger *slog.Logger
	runner pworkers.DownloadMaintenanceRunner

	// options
	enableMoveUnmatchedFiles bool
}

func NewDeleteMissingFilesJob(logger *slog.Logger, runner pworkers.DownloadMaintenanceRunner, enableMoveUnmatchedFiles bool) *deleteMissingFilesJob {
	return &deleteMissingFilesJob{
		logger: logger,
		runner: runner,

		// options
		enableMoveUnmatchedFiles: enableMoveUnmatchedFiles,
	}
}

func (j *deleteMissingFilesJob) Execute(ctx context.Context) error {
	startTime := time.Now()
	err := j.runner.DeleteMissingFiles(ctx, j.enableMoveUnmatchedFiles)
	elapsed := time.Since(startTime)

	j.logger.Debug(
		"Job done",
		"name", "DeleteMissingFiles",
		"elapsed", uformat.DurationFormat(elapsed),
	)

	return err
}

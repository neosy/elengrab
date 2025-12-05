package wjobs

import (
	"context"

	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

type deleteMissingFilesJob struct {
	runner pworkers.MaintenanceRunner

	// options
	enableMoveUnmatchedFiles bool
}

func NewDeleteMissingFilesJob(runner pworkers.MaintenanceRunner, enableMoveUnmatchedFiles bool) *deleteMissingFilesJob {
	return &deleteMissingFilesJob{
		runner: runner,

		// options
		enableMoveUnmatchedFiles: enableMoveUnmatchedFiles,
	}
}

func (j *deleteMissingFilesJob) Execute(ctx context.Context) error {
	return j.runner.DeleteMissingFiles(ctx, j.enableMoveUnmatchedFiles)
}

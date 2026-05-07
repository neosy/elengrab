package wjobs

import (
	"context"
	"log/slog"

	nworkers "github.com/neosy/elengrab/internal/pkg/workers"
	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

func NewDeleteMissingFilesJob(
	logger *slog.Logger,
	runner pworkers.DownloadMaintenanceRunner,

	// options
	enableMoveUnmatchedFiles bool,
) nworkers.Job {
	run := func(ctx context.Context) error {
		return runner.DeleteMissingFiles(ctx, enableMoveUnmatchedFiles)
	}
	return nworkers.NewJob(
		"DeleteMissingFiles",
		nworkers.MakeTimedJobExecute(logger, run),
	)
}

package wjobs

import (
	"context"
	"log/slog"

	nworkers "github.com/neosy/elengrab/internal/pkg/workers"
	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

func NewDeleteMissingDownloadsJob(
	logger *slog.Logger,
	runner pworkers.DownloadMaintenanceRunner,

	// options
	enableMoveUnmatchedDownloads bool,
) nworkers.Job {
	run := func(ctx context.Context) error {
		return runner.DeleteMissingDownloads(ctx, enableMoveUnmatchedDownloads)
	}
	return nworkers.NewJob(
		"DeleteMissingDownloads",
		nworkers.MakeTimedJobExecute(logger, run),
	)
}

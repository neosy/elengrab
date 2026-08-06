package wjobs

import (
	"context"
	"log/slog"

	"github.com/neosy/elengrab/internal/pkg/workers"
	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

func NewDeleteMissingDownloadsJob(
	logger *slog.Logger,
	runner pworkers.DownloadMaintenanceRunner,

	// options
	enableMoveUnmatchedFiles bool,
) workers.Job {
	run := func(ctx context.Context) error {
		return runner.DeleteMissingDownloads(ctx, enableMoveUnmatchedFiles)
	}
	return workers.NewJob(
		"DeleteMissingDownloads",
		workers.WrapJobExecute(logger, run),
	)
}

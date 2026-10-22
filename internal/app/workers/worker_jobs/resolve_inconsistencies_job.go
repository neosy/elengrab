package wjobs

import (
	"context"
	"log/slog"
	"time"

	"github.com/neosy/elengrab/internal/pkg/workers"
	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

func NewResolveInconsistenciesJob(
	logger *slog.Logger,
	runner pworkers.DownloadMaintenanceRunner,

	// Options
	orphanGracePeriod time.Duration,
) workers.Job {
	run := func(ctx context.Context) error {
		return runner.ResolveInconsistencies(ctx, orphanGracePeriod)
	}
	return workers.NewJob(
		"ResolveInconsistencies",
		workers.WrapJobExecute(logger, run),
	)
}

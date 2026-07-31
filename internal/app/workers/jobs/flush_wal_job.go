package wjobs

import (
	"context"
	"log/slog"

	"github.com/neosy/elengrab/internal/pkg/workers"
	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

func NewFlushWALJob(logger *slog.Logger, runner pworkers.DBMaintenanceRunner) workers.Job {
	run := func(context.Context) error {
		return runner.FlushWAL()
	}

	return workers.NewJob(
		"FlushWAL",
		workers.WrapJobExecute(logger, run),
	)
}

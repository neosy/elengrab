package wjobs

import (
	"context"
	"log/slog"

	nworkers "github.com/neosy/elengrab/internal/pkg/workers"
	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

func NewFlushWALJob(logger *slog.Logger, runner pworkers.DBMaintenanceRunner) nworkers.Job {
	run := func(context.Context) error {
		return runner.FlushWAL()
	}

	return nworkers.NewJob(
		"FlushWAL",
		nworkers.WrapJobExecute(logger, run),
	)
}

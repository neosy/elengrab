package wjobs

import (
	"log/slog"

	nworkers "github.com/neosy/elengrab/internal/pkg/workers"
	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

func NewStartupAuthWebJob(logger *slog.Logger, runner pworkers.AuthWebStartupRunner) nworkers.Job {
	return nworkers.NewJob(
		"StartupAuthWeb",
		nworkers.MakeTimedJobExecute(logger, runner.Startup),
	)
}

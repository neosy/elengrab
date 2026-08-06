package wjobs

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/pkg/workers"
	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

func NewStartupAuthWebJob(logger *slog.Logger, runner pworkers.AuthWebStartupRunner) workers.Job {
	return workers.NewJob(
		"StartupAuthWeb",
		workers.WrapJobExecute(logger, runner.Startup),
	)
}

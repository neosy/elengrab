package wjobs

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/pkg/workers"
	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

func NewStartupDatabaseJob(logger *slog.Logger, runner pworkers.DBMaintenanceRunner) workers.Job {
	return workers.NewJob(
		"StartupMaintenanceDatabase",
		workers.WrapJobExecute(logger, runner.StartupDatabase),
	)
}

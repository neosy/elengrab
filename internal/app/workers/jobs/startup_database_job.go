package wjobs

import (
	"log/slog"

	nworkers "github.com/neosy/elengrab/internal/pkg/workers"
	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

func NewStartupDatabaseJob(logger *slog.Logger, runner pworkers.DBMaintenanceRunner) nworkers.Job {
	return nworkers.NewJob(
		"StartupMaintenanceDatabase",
		nworkers.WrapJobExecute(logger, runner.StartupDatabase),
	)
}

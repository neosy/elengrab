package wjobs

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/pkg/workers"
	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

func NewDeleteDuplicatesJob(logger *slog.Logger, runner pworkers.DownloadMaintenanceRunner) workers.Job {
	return workers.NewJob(
		"DeleteDuplicates",
		workers.WrapJobExecute(logger, runner.DeleteDuplicates),
	)
}

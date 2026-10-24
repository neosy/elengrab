package wjobs

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/pkg/workers"
	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

func NewDeleteFailedMediaDownloadsJob(logger *slog.Logger, runner pworkers.DownloadMaintenanceRunner) workers.Job {
	return workers.NewJob(
		"DeleteFailedMediaDownloads",
		workers.WrapJobExecute(logger, runner.DeleteFailedMediaDownloads),
	)
}

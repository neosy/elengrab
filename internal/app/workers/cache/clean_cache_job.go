package cachejobs

import (
	"log/slog"

	nworkers "github.com/neosy/elengrab/internal/pkg/workers"
	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

func NewCleanCacheJob(logger *slog.Logger, runner pworkers.CacheRunner) nworkers.Job {
	return nworkers.NewJob(
		"CleanCacheExpired",
		nworkers.MakeTimedJobExecute(logger, runner.CleanExpired),
	)
}

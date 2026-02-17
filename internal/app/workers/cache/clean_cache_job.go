package cachejobs

import (
	"context"
	"log/slog"
	"time"

	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

type cleanCacheJob struct {
	logger *slog.Logger
	runner pworkers.CacheRunner
}

func NewCleanCacheJob(logger *slog.Logger, runner pworkers.CacheRunner) *cleanCacheJob {
	return &cleanCacheJob{
		logger: logger,
		runner: runner,
	}
}

func (j *cleanCacheJob) Execute(ctx context.Context) error {
	startTime := time.Now()
	err := j.runner.CleanExpired(ctx)
	elapsed := time.Since(startTime)

	j.logger.Debug(
		"Job done",
		"name", "CleanCacheExpired",
		"elapsed", elapsed,
	)

	return err
}

package mediawatch

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	wjobs "github.com/neosy/elengrab/internal/app/workers/jobs"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	nworkerpool "github.com/neosy/elengrab/internal/pkg/workerpool"
)

func (uc *MediaWatch) updateStats(
	ctx context.Context,
	downloadID uuid.UUID,
	mediaDuration time.Duration,
) error {
	requiredChunkCount := calcRequiredChunkCount(mediaDuration)

	views, err := uc.chunk.CountViews(ctx, downloadID, requiredChunkCount)
	if err != nil {
		uc.logger.Warn(
			"Failed to count media views",
			"downloadID", downloadID,
			"requiredChunkCount", requiredChunkCount,
			"error", err,
		)
		return err
	}

	if views == 0 {
		return nil
	}

	stat := &ddownload.MediaWatchStat{
		DownloadID: downloadID,
		Views:      views,
	}

	err = uc.stat.Write(ctx, stat)
	if err != nil {
		uc.logger.Warn(
			"Failed to update media watch statistics",
			"downloadID", downloadID,
			"views", views,
			"error", err,
		)
		return err
	}

	return nil
}

func (uc *MediaWatch) updateAllStats(ctx context.Context) error {
	var hasError bool

	if uc.pendingStats.count() == 0 {
		return nil
	}

	queue := uc.pendingStats.drain()

	for downloadID, duration := range queue {
		err := uc.updateStats(ctx, downloadID, duration)
		if err != nil {
			hasError = true
		}
	}

	if hasError {
		return errors.New("failed to update one or more media watch statistics")
	}

	return nil
}

func (uc *MediaWatch) enqueueUpdateAllStats() nworkerpool.Job {
	job := wjobs.NewUpdateWatchStatsJob(uc)

	if !uc.watchEventDispatcher.AddJob(job) {
		return nil
	}

	return job
}

func (uc *MediaWatch) startStatsUpdater(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(statsUpdateInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if uc.pendingStats.count() == 0 {
					continue
				}

				job := uc.enqueueUpdateAllStats()
				if job == nil {
					uc.logger.Warn(
						"Task has not been added to the queue",
						"name", "UpdateWatchStats",
					)
				}
			}
		}
	}()
}

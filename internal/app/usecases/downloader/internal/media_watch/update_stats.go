package mediawatch

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	wjobs "github.com/neosy/elengrab/internal/app/workers/jobs"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	"github.com/neosy/elengrab/internal/pkg/workerpool"
)

func (uc *MediaWatch) updateStats(
	ctx context.Context,
	downloadID uuid.UUID,
	mediaDuration time.Duration,
) error {
	requiredChunkCount := calcRequiredChunkCount(mediaDuration)

	views, err := uc.userChunk.CountViews(ctx, downloadID, requiredChunkCount)
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

func (uc *MediaWatch) updateUserStats(
	ctx context.Context,
	downloadID uuid.UUID,
	userID uuid.UUID,
	mediaDuration time.Duration,
) error {
	requiredChunkCount := calcRequiredChunkCount(mediaDuration)

	views, err := uc.userChunk.CountUserViews(ctx, downloadID, userID, requiredChunkCount)
	if err != nil {
		uc.logger.Warn(
			"Failed to count media views",
			"downloadID", downloadID,
			"userID", userID,
			"requiredChunkCount", requiredChunkCount,
			"error", err,
		)
		return err
	}

	if views == 0 {
		return nil
	}

	stat := &ddownload.MediaUserWatchStat{
		DownloadID: downloadID,
		UserID:     userID,
		Views:      views,
	}

	err = uc.userStat.Write(ctx, stat)
	if err != nil {
		uc.logger.Warn(
			"Failed to update media user watch statistics",
			"downloadID", downloadID,
			"userID", userID,
			"views", views,
			"error", err,
		)
		return err
	}

	return nil
}

func (uc *MediaWatch) updateStatsFromQueue(ctx context.Context) error {
	var hasError bool

	if uc.pendingStats.count() == 0 {
		return nil
	}

	// Drain the pending stats queue
	queue := uc.pendingStats.drain()

	downloadIDs := make(map[uuid.UUID]time.Duration)

	// Update user stats and collect download durations
	for key, duration := range queue {
		downloadIDs[key.downloadID] = duration
		err := uc.updateUserStats(ctx, key.downloadID, key.userID, duration)
		if err != nil {
			hasError = true
		}
	}

	// Update overall stats for each downloadID
	for downloadID, duration := range downloadIDs {
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

func (uc *MediaWatch) enqueueUpdateAllStats() workerpool.Job {
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

package mediawatch

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	wjobs "github.com/neosy/elengrab/internal/app/workers/pool_jobs"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	"github.com/neosy/elengrab/internal/pkg/workerpool"
)

func (uc *MediaWatch) updateStats(
	ctx context.Context,
	downloadID uuid.UUID,
	mediaDuration time.Duration,
) (*ddownload.MediaWatchStat, error) {
	requiredChunkCount := calcRequiredChunkCount(mediaDuration)

	views, err := uc.userChunk.CountViews(ctx, downloadID, requiredChunkCount)
	if err != nil {
		uc.logger.Warn(
			"Failed to count media views",
			"downloadID", downloadID,
			"requiredChunkCount", requiredChunkCount,
			"error", err,
		)
		return nil, err
	}

	if views == 0 {
		return nil, nil
	}

	var updatedStat *ddownload.MediaWatchStat

	err = uc.stat.Tx(
		ctx,
		func(ctx context.Context) error {
			stat, _ := uc.stat.Find(ctx, downloadID)

			if stat != nil && stat.Views == views {
				return nil
			}

			stat = &ddownload.MediaWatchStat{
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

			updatedStat, _ = uc.stat.Find(ctx, downloadID)

			return nil
		},
	)

	return updatedStat, err
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

	processed := make(map[downloadUserKey]struct{})

	// Update user stats and collect download durations
	for key, duration := range queue {
		processedKey := key.downloadUserKey()
		if _, exists := processed[processedKey]; exists {
			continue
		}
		processed[processedKey] = struct{}{}

		downloadIDs[key.downloadID] = duration
		err := uc.updateUserStats(ctx, key.downloadID, key.userID, duration)
		if err != nil {
			hasError = true
		}

	}

	updatedIDs := make(map[uuid.UUID]*ddownload.MediaWatchStat)

	// Update overall stats for each downloadID
	for downloadID, duration := range downloadIDs {
		stat, err := uc.updateStats(ctx, downloadID, duration)
		if err != nil {
			hasError = true
		}

		if stat != nil {
			updatedIDs[downloadID] = stat
		}
	}

	if hasError {
		return errors.New("failed to update one or more media watch statistics")
	}

	// Update the external handler about updated user watch statistics for each processed downloadID and userID.
	for key := range queue {
		if _, exists := updatedIDs[key.downloadID]; exists {
			uc.onWatchUserStatsUpdated(ctx, key.authCtx(), key.downloadID)
		}
	}

	// Update the external handler about updated watch statistics for each downloadID.
	for _, stat := range updatedIDs {
		uc.onWatchStatsUpdated(ctx, stat)
	}

	return nil
}

func (uc *MediaWatch) enqueueUpdateAllStats() workerpool.Job {
	job := wjobs.NewUpdateWatchStatsJob(uc)

	if err := uc.watchEventDispatcher.AddJob(job); err != nil {
		uc.logger.Error(
			"Failed to enqueue update watch stats job",
			"error", err,
		)
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

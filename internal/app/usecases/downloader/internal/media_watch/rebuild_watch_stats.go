package mediawatch

import (
	"context"
	"time"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (uc *MediaWatch) RebuildWatchStats(
	ctx context.Context,
	findDownload func(ctx context.Context, downloadID uuid.UUID) (*ddownload.MediaDownload, error),
) error {
	// Delete all existing user watch stats
	err := uc.userStat.DeleteAll(ctx)
	if err != nil {
		return err
	}

	// Delete all existing watch stats
	err = uc.stat.DeleteAll(ctx)
	if err != nil {
		return err
	}

	// Build a map of downloadID -> media duration
	downloadDurations := make(map[uuid.UUID]time.Duration)

	// Build a map witch key downloadID, userID
	downloadUsers := make(map[uuid.UUID]map[uuid.UUID]struct{})

	// Iterate through all download users and populate the maps
	err = uc.userChunk.IterateDownloadUsers(ctx,
		func(downloadID, userID uuid.UUID) error {
			var download *ddownload.MediaDownload

			if _, ok := downloadDurations[downloadID]; !ok {
				var err error
				download, err = findDownload(ctx, downloadID)
				if err != nil {
					return err
				}
				if download == nil || download.MediaInfo == nil {
					return nil
				}
				downloadDurations[downloadID] = download.MediaInfo.Duration()
			}

			if _, ok := downloadUsers[downloadID]; !ok {
				downloadUsers[downloadID] = make(map[uuid.UUID]struct{})
			}
			downloadUsers[downloadID][userID] = struct{}{}

			return nil
		},
	)
	if err != nil {
		return err
	}

	// Update stats for each downloadID and its associated users
	for downloadID, users := range downloadUsers {
		duration, ok := downloadDurations[downloadID]
		if !ok {
			continue
		}

		// Update user stats for each user associated with the downloadID
		for userID := range users {
			err := uc.updateUserStats(ctx, downloadID, userID, duration)
			if err != nil {
				return err
			}
		}

		// Update overall stats for the downloadID
		_, err = uc.updateStats(ctx, downloadID, duration)
		if err != nil {
			return err
		}
	}

	return nil
}

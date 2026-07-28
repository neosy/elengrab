package migrations

import (
	"context"
	"fmt"
	"sync"
)

var (
	rebuildWatchStatsOnceOnceSync sync.Once
)

func (m *migrations) rebuildWatchStatsOnce(ctx context.Context) (bool, error) {
	var (
		ok, executed bool
		err          error
	)

	rebuildWatchStatsOnceOnceSync.Do(func() {
		ok, err = m.rebuildWatchStats(ctx)
		executed = true
	})

	if executed {
		return ok, err
	}

	return true, nil
}

func (m *migrations) rebuildWatchStats(ctx context.Context) (bool, error) {
	m.logger.Info("Rebuilding user watch statistics...")

	err := m.usecases.mediaWatch.RebuildWatchStats(ctx, m.usecases.download.FindByDownloadID)
	if err != nil {
		m.logger.Warn(
			"Failed to rebuild user watch statistics",
			"error", err,
		)
		return false, fmt.Errorf("errors in the migration process")
	}

	m.logger.Info("User watch statistics successfully rebuilt")

	return true, nil
}

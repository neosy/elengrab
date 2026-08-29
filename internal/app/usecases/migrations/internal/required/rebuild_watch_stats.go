package required

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
	m.Logger().Info("Rebuilding user watch statistics...")

	err := m.Usecases().MediaWatch.RebuildWatchStats(ctx)
	if err != nil {
		m.Logger().Warn(
			"Failed to rebuild user watch statistics",
			"error", err,
		)
		return false, fmt.Errorf("errors in the migration process")
	}

	m.Logger().Info("User watch statistics successfully rebuilt")

	return true, nil
}

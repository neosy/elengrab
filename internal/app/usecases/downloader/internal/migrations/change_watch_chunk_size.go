package migrations

import (
	"context"
	"fmt"
	"sync"
)

var (
	changeWatchChunkSizeOnceSync sync.Once
)

func (m *migrations) changeWatchChunkSizeOnce(ctx context.Context) (bool, error) {
	var (
		ok, executed bool
		err          error
	)

	changeWatchChunkSizeOnceSync.Do(func() {
		ok, err = m.changeWatchChunkSize(ctx)
		executed = true
	})

	if executed {
		return ok, err
	}

	return true, nil
}

func (m *migrations) changeWatchChunkSize(ctx context.Context) (bool, error) {
	m.logger.Info("Changing media watch chunk size...")

	err := m.usecases.mediaWatch.RebuildUserChunks(ctx, m.usecases.download.FindByDownloadID)
	if err != nil {
		m.logger.Warn(
			"Failed to rebuild media watch chunks",
			"error", err,
		)
		return false, fmt.Errorf("errors in the migration process")
	}

	m.logger.Info("Media watch chunk size successfully changed")

	return true, nil
}

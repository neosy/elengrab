package required

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
	m.Logger().Info("Changing media watch chunk size...")

	err := m.Usecases().MediaWatch.RebuildUserChunks(ctx)
	if err != nil {
		m.Logger().Warn(
			"Failed to rebuild media watch chunks",
			"error", err,
		)
		return false, fmt.Errorf("errors in the migration process")
	}

	m.Logger().Info("Media watch chunk size successfully changed")

	return true, nil
}

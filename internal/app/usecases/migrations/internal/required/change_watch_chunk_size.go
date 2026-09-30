package required

import (
	"context"
	"fmt"
	"sync"
)

var (
	changeWatchChunkSizeOnceSync sync.Once
)

func (m *migrations) changeWatchChunkSizeOnce(ctx context.Context) error {
	var err error

	changeWatchChunkSizeOnceSync.Do(func() {
		err = m.changeWatchChunkSize(ctx)
	})

	return err
}

func (m *migrations) changeWatchChunkSize(ctx context.Context) error {
	m.Logger().Info("Changing media watch chunk size...")

	err := m.Usecases().MediaWatch.RebuildUserChunks(ctx)
	if err != nil {
		m.Logger().Warn(
			"Failed to rebuild media watch chunks",
			"error", err,
		)
		return fmt.Errorf("errors in the migration process")
	}

	m.Logger().Info("Media watch chunk size successfully changed")

	return nil
}

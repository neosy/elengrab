package migrations

import (
	"context"
	"fmt"
	"time"

	mediawatch "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/media_watch"
)

func (m *migrations) changeWatchChunkSizeTo500ms(ctx context.Context) (bool, error) {
	newSize := 500 * time.Millisecond

	if newSize != mediawatch.ChunkDuration {
		return true, nil
	}

	return m.changeWatchChunkSize(ctx)
}

func (m *migrations) changeWatchChunkSize(ctx context.Context) (bool, error) {
	m.logger.Info("Changing media watch chunk size...")

	err := m.usecases.mediaWatch.RebuildChunks(ctx, m.usecases.download.FindByDownloadID)
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

package mediawatch

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (uc *MediaWatch) RebuildUserChunks(
	ctx context.Context,
	findDownload func(ctx context.Context, downloadID uuid.UUID) (*ddownload.MediaDownload, error),
) error {
	err := uc.userChunk.DeleteAll(ctx)
	if err != nil {
		return err
	}

	chunkBatches := make([][]*ddownload.MediaUserWatchChunk, 0)

	err = uc.event.IterateGetAll(
		ctx,
		func(event *ddownload.MediaWatchEvent) error {
			download, err := findDownload(ctx, event.DownloadID)
			if err != nil {
				return err
			}

			if download == nil || download.MediaInfo == nil {
				return nil
			}

			chunks := uc.eventToChunks(event, download.MediaInfo.Duration())
			if len(chunks) == 0 {
				return nil
			}

			chunkBatches = append(chunkBatches, chunks)

			return nil
		},
	)
	if err != nil {
		return err
	}

	for _, chunks := range chunkBatches {
		err = uc.userChunk.AddChunkQtyBatch(ctx, chunks)
		if err != nil {
			return err
		}
	}

	return nil
}

package uwatchchunk

import (
	"context"

	"uuid"
)

func (uc *MediaUserWatchChunk) DeleteByDownloadID(ctx context.Context, downloadID uuid.UUID) error {
	return uc.chunkRepo().DeleteByDownloadID(ctx, downloadID)
}

func (uc *MediaUserWatchChunk) DeleteAll(ctx context.Context) error {
	return uc.chunkRepo().DeleteAll(ctx)
}

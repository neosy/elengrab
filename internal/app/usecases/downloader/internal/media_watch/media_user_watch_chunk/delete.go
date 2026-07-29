package uwatchchunk

import (
	"context"

	"github.com/google/uuid"
)

func (uc *MediaUserWatchChunk) DeleteByDownloadID(ctx context.Context, downloadID uuid.UUID) error {
	return uc.chunkRep.DeleteByDownloadID(ctx, downloadID)
}

func (uc *MediaUserWatchChunk) DeleteAll(ctx context.Context) error {
	return uc.chunkRep.DeleteAll(ctx)
}

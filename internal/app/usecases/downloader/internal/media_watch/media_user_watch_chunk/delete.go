package uwatchchunk

import (
	"context"

	"github.com/google/uuid"
)

func (uc *MediaUserWatchChunk) Delete(ctx context.Context, downloadID uuid.UUID) error {
	return uc.chunkRep.Delete(ctx, downloadID)
}

func (uc *MediaUserWatchChunk) DeleteAll(ctx context.Context) error {
	return uc.chunkRep.DeleteAll(ctx)
}

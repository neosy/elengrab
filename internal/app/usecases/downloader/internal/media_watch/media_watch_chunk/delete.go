package watchchunk

import (
	"context"

	"github.com/google/uuid"
)

func (uc *MediaWatchChunk) Delete(ctx context.Context, downloadID uuid.UUID) error {
	return uc.chunkRep.Delete(ctx, downloadID)
}

func (uc *MediaWatchChunk) DeleteAll(ctx context.Context) error {
	return uc.chunkRep.DeleteAll(ctx)
}

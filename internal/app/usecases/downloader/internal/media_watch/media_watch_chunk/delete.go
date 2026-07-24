package watchchunk

import (
	"context"

	"github.com/google/uuid"
)

func (uc *MediaWatchChunk) Delete(ctx context.Context, downloadID uuid.UUID) error {
	return uc.chunkRep.Delete(ctx, downloadID)
}

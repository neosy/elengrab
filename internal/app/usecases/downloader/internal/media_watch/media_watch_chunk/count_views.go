package watchchunk

import (
	"context"

	"github.com/google/uuid"
)

func (uc *MediaWatchChunk) CountViews(ctx context.Context, downloadID uuid.UUID, requiredChunks uint32) (uint32, error) {
	return uc.chunkRep.CountViews(ctx, downloadID, requiredChunks)
}

func (uc *MediaWatchChunk) CountUserViews(
	ctx context.Context,
	downloadID uuid.UUID, userID uuid.UUID,
	requiredChunks uint32) (uint32, error) {
	return uc.chunkRep.CountUserViews(ctx, downloadID, userID, requiredChunks)
}

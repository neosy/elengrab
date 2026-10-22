package uwatchchunk

import (
	"context"

	"uuid"
)

func (uc *MediaUserWatchChunk) CountViews(ctx context.Context, downloadID uuid.UUID, requiredChunks uint32) (uint32, error) {
	return uc.chunkRepo().CountViews(ctx, downloadID, requiredChunks)
}

func (uc *MediaUserWatchChunk) CountUserViews(
	ctx context.Context,
	downloadID uuid.UUID, userID uuid.UUID,
	requiredChunks uint32) (uint32, error) {
	return uc.chunkRepo().CountUserViews(ctx, downloadID, userID, requiredChunks)
}

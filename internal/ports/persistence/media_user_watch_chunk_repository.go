package persistence

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

type MediaUserWatchChunkRepository interface {
	Transactional

	AddChunkQty(ctx context.Context, chunk *ddownload.MediaUserWatchChunk) error
	AddChunkQtyBatch(ctx context.Context, chunks []*ddownload.MediaUserWatchChunk) error
	Delete(ctx context.Context, downloadID uuid.UUID) error
	DeleteAll(ctx context.Context) error

	IterateDownloadUsers(ctx context.Context, fn func(downloadID, userID uuid.UUID) error) error

	CountViews(ctx context.Context, downloadID uuid.UUID, requiredChunks uint32) (uint32, error)
	CountUserViews(ctx context.Context, downloadID uuid.UUID, userID uuid.UUID, requiredChunks uint32) (uint32, error)

	WithUserID() MediaUserWatchChunkRepository
}

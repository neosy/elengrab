package persistence

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

type MediaWatchChunkRepository interface {
	Transactional

	AddChunkQty(ctx context.Context, chunk *ddownload.MediaWatchChunk) error
	AddChunkQtyBatch(ctx context.Context, chunks []*ddownload.MediaWatchChunk) error
	Delete(ctx context.Context, downloadID uuid.UUID) error
	DeleteAll(ctx context.Context) error

	CountViews(ctx context.Context, downloadID uuid.UUID, requiredChunks uint32) (uint32, error)

	WithUserID() MediaWatchChunkRepository
}

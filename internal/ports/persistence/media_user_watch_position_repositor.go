package persistence

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	"github.com/neosy/elengrab/internal/pkg/cache/memory"
	memsimple "github.com/neosy/elengrab/internal/pkg/cache/memory/simple"
)

type MediaUserWatchPositionRepository interface {
	Transactional

	// Insert inserting a record
	Write(ctx context.Context, position *ddownload.MediaUserWatchPosition) error
	DeleteByDownloadID(ctx context.Context, downloadID uuid.UUID) error

	Find(ctx context.Context, downloadID uuid.UUID, userID uuid.UUID, sessionID uuid.UUID) (*ddownload.MediaUserWatchPosition, error)
}

type MediaUserWatchPositionCacheRepository interface {
	memory.CacheRepository

	Save(ctx context.Context, position *ddownload.MediaUserWatchPosition) error
	SaveNegative(ctx context.Context, downloadID uuid.UUID, userID uuid.UUID, sessionID uuid.UUID) error
	Delete(ctx context.Context, downloadID uuid.UUID, userID uuid.UUID, sessionID uuid.UUID) error

	Find(ctx context.Context, downloadID uuid.UUID, userID uuid.UUID, sessionID uuid.UUID) (*ddownload.MediaUserWatchPosition, memsimple.CacheStatus, error)
	Exists(ctx context.Context, downloadID uuid.UUID, userID uuid.UUID, sessionID uuid.UUID) (bool, error)

	CleanExpired(context.Context) error
}

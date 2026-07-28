package persistence

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	memsimple "github.com/neosy/elengrab/internal/pkg/cache/memory/simple"
)

type MediaUserWatchPositionRepository interface {
	Transactional

	// Insert inserting a record
	Write(ctx context.Context, position *ddownload.MediaUserWatchPosition) error
	DeleteByDownloadID(ctx context.Context, downloadID uuid.UUID) error

	Find(ctx context.Context, downloadID uuid.UUID, userID uuid.UUID, sessionID *uuid.UUID) (*ddownload.MediaUserWatchPosition, error)
}

type MediaUserWatchPositionCacheRepository interface {
	Save(position *ddownload.MediaUserWatchPosition) error
	SaveNegative(downloadID uuid.UUID, userID uuid.UUID, sessionID *uuid.UUID) error
	Delete(downloadID uuid.UUID, userID uuid.UUID, sessionID *uuid.UUID) error

	Find(downloadID uuid.UUID, userID uuid.UUID, sessionID *uuid.UUID) (*ddownload.MediaUserWatchPosition, memsimple.CacheStatus, error)
	Exists(downloadID uuid.UUID, userID uuid.UUID, sessionID *uuid.UUID) (bool, error)

	CleanExpired(context.Context) error
}

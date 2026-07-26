package persistence

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	memsimple "github.com/neosy/elengrab/internal/pkg/cache/memory/simple"
)

type MediaWatchPositionRepository interface {
	Transactional

	// Insert inserting a record
	Write(ctx context.Context, position *ddownload.MediaWatchPosition) error
	DeleteByDownloadID(ctx context.Context, downloadID uuid.UUID) error

	Find(ctx context.Context, downloadID uuid.UUID, userID uuid.UUID, sessionID *uuid.UUID) (*ddownload.MediaWatchPosition, error)
}

type MediaWatchPositionCacheRepository interface {
	Save(position *ddownload.MediaWatchPosition) error
	SaveNegative(downloadID uuid.UUID, userID uuid.UUID, sessionID *uuid.UUID) error
	Delete(downloadID uuid.UUID, userID uuid.UUID, sessionID *uuid.UUID) error

	Find(downloadID uuid.UUID, userID uuid.UUID, sessionID *uuid.UUID) (*ddownload.MediaWatchPosition, memsimple.CacheStatus, error)
	Exists(downloadID uuid.UUID, userID uuid.UUID, sessionID *uuid.UUID) (bool, error)

	CleanExpired(context.Context) error
}

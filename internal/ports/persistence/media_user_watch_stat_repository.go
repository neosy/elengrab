package persistence

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	"github.com/neosy/elengrab/internal/pkg/cache/memory"
	memsimple "github.com/neosy/elengrab/internal/pkg/cache/memory/simple"
)

type MediaUserWatchStatRepository interface {
	Transactional

	// Insert inserting a record
	Write(ctx context.Context, stat *ddownload.MediaUserWatchStat) error
	Delete(ctx context.Context, downloadID uuid.UUID, userID uuid.UUID) error
	DeleteByDownloadID(ctx context.Context, downloadID uuid.UUID) error
	DeleteAll(ctx context.Context) error

	Find(ctx context.Context, downloadID uuid.UUID, userID uuid.UUID) (*ddownload.MediaUserWatchStat, error)
	Exists(ctx context.Context, downloadID uuid.UUID, userID uuid.UUID) (bool, error)
}

type MediaUserWatchStatCacheRepository interface {
	memory.CacheRepository

	Save(stat *ddownload.MediaUserWatchStat) error
	SaveNegative(downloadID, userID uuid.UUID) error
	Delete(downloadID, userID uuid.UUID) error

	Find(downloadID, userID uuid.UUID) (*ddownload.MediaUserWatchStat, memsimple.CacheStatus, error)
	Exists(downloadID, userID uuid.UUID) (bool, memsimple.CacheStatus, error)

	CleanExpired(context.Context) error
}

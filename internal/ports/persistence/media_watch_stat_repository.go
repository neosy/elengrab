package persistence

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	memsimple "github.com/neosy/elengrab/internal/pkg/cache/memory/simple"
)

type MediaWatchStatRepository interface {
	Transactional

	// Insert inserting a record
	Write(ctx context.Context, stat *ddownload.MediaWatchStat) error
	DeleteByDownloadID(ctx context.Context, downloadID uuid.UUID) error
	DeleteAll(ctx context.Context) error

	Find(ctx context.Context, downloadID uuid.UUID) (*ddownload.MediaWatchStat, error)
}

type MediaWatchStatCacheRepository interface {
	Save(stat *ddownload.MediaWatchStat) error
	SaveNegative(downloadID uuid.UUID) error
	Delete(downloadID uuid.UUID) error

	Find(downloadID uuid.UUID) (*ddownload.MediaWatchStat, memsimple.CacheStatus, error)
	Exists(downloadID uuid.UUID) (bool, error)

	CleanExpired(context.Context) error
}

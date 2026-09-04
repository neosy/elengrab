package persistence

import (
	"context"
	"time"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/cache/memory"
	memsimple "github.com/neosy/elengrab/internal/pkg/cache/memory/simple"
)

type MediaDownloadRepositoryFactory func() MediaDownloadRepository

type MediaDownloadRepository interface {
	Transactional
	Insert(ctx context.Context, download *ddownload.MediaDownload) error
	Update(ctx context.Context, download *ddownload.MediaDownload) error
	SoftDelete(ctx context.Context, DownloadID uuid.UUID) error
	HardDelete(ctx context.Context, DownloadID uuid.UUID) error
	Restore(ctx context.Context, DownloadID uuid.UUID) error

	// UpdateStatus updates all jobs with status [Working or Pending to New], [Refreshing to Donw].
	UpdateStatus(
		ctx context.Context,
		fromStatuses []dtypes.MediaDownloadStatus,
		newStatus dtypes.MediaDownloadStatus,
	) error
	FillEmptyMediaTitleLower(ctx context.Context) error
	UpdateOwner(ctx context.Context, fromID, toID uuid.UUID) error

	FindByDownloadID(ctx context.Context, DownloadID uuid.UUID) (*ddownload.MediaDownload, error)
	IterateGetAll(ctx context.Context, fn func(*ddownload.MediaDownload) error) error
	GetAllFullNames(ctx context.Context, includeDeleted bool) (map[string]struct{}, error)
	IterateFullNames(ctx context.Context, includeDeleted bool, fn func(string) error) error
	IterateGetByIDs(ctx context.Context, ids []uuid.UUID, fn func(*ddownload.MediaDownload) error) error
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*ddownload.MediaDownload, error)
	GetByStatus(ctx context.Context, status dtypes.MediaDownloadStatus) ([]*ddownload.MediaDownload, error)
	GetByStatuses(ctx context.Context, statuses []dtypes.MediaDownloadStatus) ([]*ddownload.MediaDownload, error)
	GetByPartialHash(ctx context.Context, hash string) ([]*ddownload.MediaDownload, error)
	GetWithoutPartialHash(ctx context.Context) ([]*ddownload.MediaDownload, error)
	GetDuplicateHashes(ctx context.Context, scope dtypes.UniquenessScope) ([]ddownload.DuplicateHashRow, error)
	GetDeleted(ctx context.Context, from, to *time.Time) ([]*ddownload.MediaDownload, error)

	WithOptions(options dtypes.QueryMediaOptions) MediaDownloadRepository
	WithStatus(statuses ...dtypes.MediaDownloadStatus) MediaDownloadRepository
	WithUser(userID uuid.UUID) MediaDownloadRepository
	WithDeleted() MediaDownloadRepository
	WithFilters(filters map[dtypes.QueryFilterName]any) MediaDownloadRepository
}

type MediaDownloadCacheRepository interface {
	memory.CacheRepository

	Save(ctx context.Context, media *ddownload.MediaDownload) error
	SaveNegative(ctx context.Context, downloadID uuid.UUID) error
	Delete(ctx context.Context, downloadID uuid.UUID) error

	FindByDownloadID(ctx context.Context, downloadID uuid.UUID) (*ddownload.MediaDownload, memsimple.CacheStatus, error)
	ExistsByFileID(ctx context.Context, downloadID uuid.UUID) (bool, error)

	CleanExpired(context.Context) error
}

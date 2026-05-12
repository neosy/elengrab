package persistence

import (
	"context"
	"time"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type MediaDownloadRepository interface {
	Transactional
	Insert(ctx context.Context, download *ddownload.MediaDownload) error
	Update(ctx context.Context, download *ddownload.MediaDownload) error
	Delete(ctx context.Context, DownloadID uuid.UUID, soft bool) error
	Restore(ctx context.Context, DownloadID uuid.UUID) error

	// UpdateStatusToNew updates all jobs with status Working or Pending to New.
	UpdateStatusToNew(ctx context.Context, statuses []dtypes.MediaDownloadStatus) error
	FillEmptyMediaTitleLower(ctx context.Context) error
	UpdateOwner(ctx context.Context, fromID, toID uuid.UUID) error

	FindByDownloadID(ctx context.Context, DownloadID uuid.UUID) (*ddownload.MediaDownload, error)
	IterateGetAll(ctx context.Context, includeDeleted bool, fn func(*ddownload.MediaDownload) error) error
	GetAllFullNames(ctx context.Context, includeDeleted bool) (map[string]struct{}, error)
	IterateFullNames(ctx context.Context, includeDeleted bool, fn func(string) error) error
	GetBeforeTime(ctx context.Context, before time.Time, limit uint64) ([]*ddownload.MediaDownload, error)
	GetByStatus(ctx context.Context, status dtypes.MediaDownloadStatus) ([]*ddownload.MediaDownload, error)
	GetByStatuses(ctx context.Context, statuses []dtypes.MediaDownloadStatus) ([]*ddownload.MediaDownload, error)
	GetByPartialHash(ctx context.Context, hash string) ([]*ddownload.MediaDownload, error)
	GetWithoutPartialHash(ctx context.Context) ([]*ddownload.MediaDownload, error)
	GetDuplicateHashes(ctx context.Context, scope dtypes.UniquenessScope) ([]ddownload.DuplicateHashRow, error)
	GetDeleted(ctx context.Context, from, to *time.Time) ([]*ddownload.MediaDownload, error)

	WithUser(userID uuid.UUID) MediaDownloadRepository
	WithFilters(filters map[string]any) MediaDownloadRepository
}

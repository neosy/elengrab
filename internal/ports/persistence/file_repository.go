package persistence

import (
	"context"
	"time"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type FileRepository interface {
	Transactional

	Insert(ctx context.Context, file *ddownload.File) error
	Update(ctx context.Context, file *ddownload.File) error
	// UpdateStatusToNew updates all jobs with status Working or Pending to New.
	UpdateStatusToNew(ctx context.Context, statuses []dtypes.FileStatus) error
	Delete(ctx context.Context, fileId uuid.UUID, soft bool) error
	Restore(ctx context.Context, fileId uuid.UUID) error
	FindByFileId(ctx context.Context, fileId uuid.UUID) (*ddownload.File, error)
	GetAll(ctx context.Context, includeDeleted bool) ([]*ddownload.File, error)
	GetAllFullNames(ctx context.Context, includeDeleted bool) ([]string, error)
	GetBeforeTime(ctx context.Context, before time.Time, limit uint64) ([]*ddownload.File, error)
	GetByStatus(ctx context.Context, status dtypes.FileStatus) ([]*ddownload.File, error)
	GetByStatuses(ctx context.Context, statuses []dtypes.FileStatus) ([]*ddownload.File, error)
	GetByPartialHash(ctx context.Context, hash string) ([]*ddownload.File, error)
	GetWithoutPartialHash(ctx context.Context) ([]*ddownload.File, error)
	GetDuplicateHashes(ctx context.Context) ([]string, error)
	GetDeleted(ctx context.Context, from, to *time.Time) ([]*ddownload.File, error)
}

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
	UpdateStatusToNew(ctx context.Context) error
	Delete(ctx context.Context, fileId uuid.UUID) error
	FindByFileId(ctx context.Context, fileId uuid.UUID) (*ddownload.File, error)
	GetAll(ctx context.Context) ([]*ddownload.File, error)
	GetBeforeTime(ctx context.Context, before time.Time, limit uint64) ([]*ddownload.File, error)
	GetByStatus(ctx context.Context, status dtypes.FileStatus) ([]*ddownload.File, error)
	GetByPartialHash(ctx context.Context, hash string) ([]*ddownload.File, error)
	GetWithoutPartialHash(ctx context.Context) ([]*ddownload.File, error)
	GetDuplicateHashes(ctx context.Context) ([]string, error)
}

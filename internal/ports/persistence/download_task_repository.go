package persistence

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type DownloadTaskRepository interface {
	Insert(ctx context.Context, task *ddownload.DownloadTask) error
	Update(ctx context.Context, task *ddownload.DownloadTask) error
	// UpdateStatusToNew updates all jobs with status Working or Pending to New.
	UpdateStatusToNew(ctx context.Context) error
	Delete(ctx context.Context, taskId uuid.UUID) error
	DeleteByFileId(ctx context.Context, fileId uuid.UUID) error
	DeleteByStatus(ctx context.Context, status dtypes.DownloadTaskStatus) error
	FindByTaskId(ctx context.Context, taskId uuid.UUID) (*ddownload.DownloadTask, error)
	FindByFileId(ctx context.Context, fileId uuid.UUID) (*ddownload.DownloadTask, error)
}

package persistence

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

type DownloadTaskRepository interface {
	Insert(ctx context.Context, task *ddownload.DownloadTask) error
	Update(ctx context.Context, task *ddownload.DownloadTask) error
	Delete(ctx context.Context, taskId uuid.UUID) error
	DeleteByFileId(ctx context.Context, fileId uuid.UUID) error
	FindByTaskId(ctx context.Context, taskId uuid.UUID) (*ddownload.DownloadTask, error)
	FindByFileId(ctx context.Context, fileId uuid.UUID) (*ddownload.DownloadTask, error)
}

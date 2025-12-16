package persistence

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

type DownloadStateRepository interface {
	Save(ctx context.Context, state *ddownload.DownloadState) error
	Insert(ctx context.Context, state *ddownload.DownloadState) error
	Update(ctx context.Context, state *ddownload.DownloadState) error
	Delete(ctx context.Context, fileId uuid.UUID) error
	FindByFileId(ctx context.Context, fileId uuid.UUID) (*ddownload.DownloadState, error)
	FindByTaskId(ctx context.Context, taskId uuid.UUID) (*ddownload.DownloadState, error)
}

type DownloadStateCacheRepository interface {
	DownloadStateRepository
	CleanExpired(ctx context.Context) error
}

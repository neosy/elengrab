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
	Delete(ctx context.Context, fileID uuid.UUID) error
	FindByFileID(ctx context.Context, fileID uuid.UUID) (*ddownload.DownloadState, error)
	FindByTaskID(ctx context.Context, taskID uuid.UUID) (*ddownload.DownloadState, error)

	WithUser(userID uuid.UUID) DownloadStateRepository
}

type DownloadStateCacheRepository interface {
	DownloadStateRepository
	CleanExpired(ctx context.Context) error
}

package persistence

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

type DownloadStateCacheRepository interface {
	Save(ctx context.Context, state *ddownload.DownloadState) error
	Insert(ctx context.Context, state *ddownload.DownloadState) error
	Update(ctx context.Context, state *ddownload.DownloadState) error
	Delete(ctx context.Context, downloadID uuid.UUID) error
	FindByDownloadID(ctx context.Context, downloadID uuid.UUID) (*ddownload.DownloadState, error)
	FindByTaskID(ctx context.Context, taskID uuid.UUID) (*ddownload.DownloadState, error)

	WithUser(userID uuid.UUID) DownloadStateCacheRepository

	CleanExpired(ctx context.Context) error
}

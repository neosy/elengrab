package persistence

import (
	"context"

	"uuid"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/cache/memory"
)

type DownloadStateCacheRepository interface {
	memory.CacheRepository

	Save(ctx context.Context, state *ddownload.DownloadState) error
	Insert(ctx context.Context, state *ddownload.DownloadState) error
	Update(ctx context.Context, state *ddownload.DownloadState) error
	Delete(ctx context.Context, downloadID uuid.UUID) error

	FindByDownloadID(ctx context.Context, downloadID uuid.UUID) (*ddownload.DownloadState, error)
	FindByTaskID(ctx context.Context, taskID uuid.UUID) (*ddownload.DownloadState, error)
	Count(ctx context.Context, userID uuid.UUID) (dtypes.IncompleteDownloadsCount, error)

	WithUser(userID uuid.UUID) DownloadStateCacheRepository

	Clear(ctx context.Context) error
}

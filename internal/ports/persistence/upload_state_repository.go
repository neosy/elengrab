package persistence

import (
	"context"

	"uuid"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/cache/memory"
)

type UploadStateCacheRepository interface {
	memory.CacheRepository

	Save(ctx context.Context, state *ddownload.UploadState) error
	Delete(ctx context.Context, downloadID uuid.UUID) error

	FindByDownloadID(ctx context.Context, downloadID uuid.UUID) (*ddownload.UploadState, error)
	Count(ctx context.Context, userID uuid.UUID) (dtypes.IncompleteDownloadsCount, error)

	WithUser(userID uuid.UUID) UploadStateCacheRepository

	Clear(ctx context.Context) error
}

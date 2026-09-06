package persistence

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type MediaSourceIndexRepositoryFactory func() MediaSourceIndexRepository

type MediaSourceIndexRepository interface {
	Transactional

	Insert(ctx context.Context, index *ddownload.MediaSourceIndex) error
	Update(ctx context.Context, index *ddownload.MediaSourceIndex) error
	Save(ctx context.Context, index *ddownload.MediaSourceIndex) error

	UpdateOwner(ctx context.Context, fromID, toID uuid.UUID) error

	SoftDelete(ctx context.Context, downloadID uuid.UUID) error
	HardDelete(ctx context.Context, downloadID uuid.UUID) error
	Restore(ctx context.Context, downloadID uuid.UUID) error

	FindByDownloadID(ctx context.Context, downloadID uuid.UUID) (*ddownload.MediaSourceIndex, error)
	IterateGetAll(ctx context.Context, fn func(*ddownload.MediaSourceIndex) error) error

	WithOptions(options dtypes.QueryMediaOptions) MediaSourceIndexRepository
	WithDeleted() MediaSourceIndexRepository
	WithUser(userID uuid.UUID) MediaSourceIndexRepository
	WithFilters(filters ...dtypes.QueryFilter) MediaSourceIndexRepository
	WithOrderBy(orderBys ...dtypes.QueryOrderBy) MediaSourceIndexRepository
}

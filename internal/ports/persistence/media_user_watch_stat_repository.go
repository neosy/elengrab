package persistence

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

type MediaUserWatchStatRepository interface {
	Transactional

	// Insert inserting a record
	Write(ctx context.Context, stat *ddownload.MediaUserWatchStat) error
	Delete(ctx context.Context, downloadID uuid.UUID, userID uuid.UUID) error
	DeleteByDownloadID(ctx context.Context, downloadID uuid.UUID) error
	DeleteAll(ctx context.Context) error

	Find(ctx context.Context, downloadID uuid.UUID, userID uuid.UUID) (*ddownload.MediaUserWatchStat, error)
}

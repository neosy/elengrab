package persistence

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

type MediaWatchEventRepository interface {
	Transactional

	// Insert inserting a record
	Insert(ctx context.Context, event *ddownload.MediaWatchEvent) error
	Update(ctx context.Context, event *ddownload.MediaWatchEvent) error
	Write(ctx context.Context, event *ddownload.MediaWatchEvent) error
	Delete(ctx context.Context, downloadID uuid.UUID) error

	IterateGetAll(ctx context.Context, fn func(*ddownload.MediaWatchEvent) error) error

	WithUserID() MediaWatchEventRepository
}

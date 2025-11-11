package persistence

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

type FileRepository interface {
	Insert(ctx context.Context, file *ddownload.File) error
	Update(ctx context.Context, file *ddownload.File) error
	Delete(ctx context.Context, fileId uuid.UUID) error
	FindByFileId(ctx context.Context, fileId uuid.UUID) (*ddownload.File, error)
}

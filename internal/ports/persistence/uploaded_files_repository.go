package persistence

import (
	"context"

	"uuid"

	"github.com/neosy/elengrab/internal/pkg/cache/memory"
	memsimple "github.com/neosy/elengrab/internal/pkg/cache/memory/simple"
	pservices "github.com/neosy/elengrab/internal/ports/services"
)

type UploadedFileCacheRepository interface {
	memory.CacheRepository

	Save(ctx context.Context, upload pservices.UploadedFile) error
	SaveNegative(ctx context.Context, fileID uuid.UUID) error

	// Delete removes the entry from the cache and deletes its file from the temporary directory.
	Delete(ctx context.Context, fileID uuid.UUID) error

	// Forget removes the entry from the cache but keeps the file on disk.
	// Use it when the file has been handed over to another owner (for example, queued for processing).
	Forget(ctx context.Context, fileID uuid.UUID) error

	Find(
		ctx context.Context,
		fileID uuid.UUID,
	) (pservices.UploadedFile, memsimple.CacheStatus, error)
	Exists(ctx context.Context, fileID uuid.UUID) (bool, error)

	CleanExpired(ctx context.Context) error
}

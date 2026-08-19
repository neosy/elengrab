package inmemory

import (
	"context"
	"time"

	"github.com/google/uuid"
	memsimple "github.com/neosy/elengrab/internal/pkg/cache/memory/simple"
)

type thumbnailRawFile struct {
	Data []byte
}

func (f *thumbnailRawFile) Copy() *thumbnailRawFile {
	if f == nil {
		return nil
	}

	file := new(*f)

	if f.Data != nil {
		file.Data = make([]byte, len(f.Data))
		copy(file.Data, f.Data)
	}

	return file
}

// Defines the structure for the in-memory repository.
type ThumbnailFileRepository struct {
	// Embeds the base Repository
	memsimple.Repository[thumbnailRawFile]

	// Cache for storing Thumbnails by their unique FileName.
	cacheByFileID memsimple.Cache[uuid.UUID, thumbnailRawFile]
}

// newThumbnailFileRepository returns a new object for the repository
func newThumbnailFileRepository(ttl time.Duration) *ThumbnailFileRepository {
	r := &ThumbnailFileRepository{
		cacheByFileID: memsimple.NewCacheWithDeaultCopier[uuid.UUID, thumbnailRawFile, *thumbnailRawFile](),
	}
	r.Repository.Init(ttl)
	return r
}

func (r *ThumbnailFileRepository) Name() string {
	return "thumbnail_file"
}

func (r *ThumbnailFileRepository) Save(ctx context.Context, fileID uuid.UUID, raw []byte) error {
	if len(raw) == 0 || fileID == uuid.Nil {
		return nil
	}

	data := make([]byte, len(raw))
	copy(data, raw)

	file := &thumbnailRawFile{
		Data: data,
	}

	save := func() error {
		r.cacheByFileID.Save(
			fileID,
			file,
			r.TTL(),
		)
		return nil
	}

	return r.Repository.Save(ctx, save)
}

func (r *ThumbnailFileRepository) SaveNegative(ctx context.Context, fileID uuid.UUID) error {
	if fileID == uuid.Nil {
		return nil
	}

	save := func() error {
		r.cacheByFileID.Save(
			fileID,
			nil,
			r.TTL(),
		)
		return nil
	}

	return r.Repository.Save(ctx, save)
}

// Delete removes a thumbnail from the in-memory repository using its ID.
func (r *ThumbnailFileRepository) Delete(ctx context.Context, fileID uuid.UUID) error {
	delete := func() error {
		if fileID != uuid.Nil {
			r.cacheByFileID.Delete(fileID)
		}
		return nil
	}
	return r.Repository.Delete(ctx, delete)
}

// FindByFileID retrieves a thumbnail by its thumbnail ID from the repository.
func (r *ThumbnailFileRepository) FindByFileID(
	ctx context.Context, fileID uuid.UUID,
) ([]byte, memsimple.CacheStatus, error) {
	find := func() (*thumbnailRawFile, memsimple.CacheStatus, error) {
		data, status := r.cacheByFileID.FindWithStatus(fileID)
		return data, status, nil
	}

	var data []byte
	file, status, err := r.Repository.FindWithStatus(ctx, find)
	if file != nil {
		data = file.Data
	}

	return data, status, err
}

// Checks if a thumbnail exists by its thumbnail ID.
func (r *ThumbnailFileRepository) ExistsByFileID(ctx context.Context, fileID uuid.UUID) (bool, error) {
	exists := func() (bool, error) {
		return r.cacheByFileID.Exists(fileID), nil
	}
	return r.Repository.Exists(ctx, exists)
}

// CleanExpired cleans expired entries from the repository.
func (r *ThumbnailFileRepository) CleanExpired(ctx context.Context) error {
	// Define a clean function to remove expired entries from the cache.
	clean := func() error {
		r.cacheByFileID.CleanExpired()
		return nil
	}
	// Call the base repository's CleanExpired method with the custom clean function.
	return r.Repository.CleanExpired(ctx, clean)
}

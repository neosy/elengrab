package inmemory

import (
	"context"
	"time"

	"uuid"

	memsimple "github.com/neosy/elengrab/internal/pkg/cache/memory/simple"
	pservices "github.com/neosy/elengrab/internal/ports/services"
)

type uploadedFile struct {
	upload pservices.UploadedFile
}

func (f *uploadedFile) Clone() *uploadedFile {
	if f == nil {
		return nil
	}

	cp := *f

	return &cp
}

type UploadedFileRepository struct {
	memsimple.Repository[uploadedFile]

	uploadedFiles memsimple.Cache[uuid.UUID, uploadedFile]
}

// newUploadedFileRepository returns a new object for the repository
func newUploadedFileRepository(ttl time.Duration) *UploadedFileRepository {
	r := &UploadedFileRepository{
		uploadedFiles: memsimple.NewCacheWithDeaultCloner[uuid.UUID, uploadedFile, *uploadedFile](),
	}
	r.Repository.Init(ttl)
	return r
}

func (r *UploadedFileRepository) Name() string {
	return "uploaded_files"
}

func (r *UploadedFileRepository) Save(ctx context.Context, upload pservices.UploadedFile) error {
	if upload == nil {
		return nil
	}

	uploadedFile := &uploadedFile{
		upload: upload,
	}

	save := func() error {
		r.uploadedFiles.Save(
			upload.FileID(),
			uploadedFile,
			r.TTL(),
		)
		return nil
	}

	return r.Repository.Save(ctx, save)
}

func (r *UploadedFileRepository) SaveNegative(ctx context.Context, fileID uuid.UUID) error {
	if fileID == uuid.Nil() {
		return nil
	}

	save := func() error {
		r.uploadedFiles.Save(
			fileID,
			nil,
			r.TTL(),
		)
		return nil
	}

	return r.Repository.Save(ctx, save)
}

// Delete removes the entry from the cache and deletes its file from the temporary directory.
func (r *UploadedFileRepository) Delete(ctx context.Context, fileID uuid.UUID) error {
	return r.remove(ctx, fileID, true)
}

// Forget removes the entry from the cache but keeps the file on disk.
// Use it when the file has been handed over to another owner (for example, queued for processing).
func (r *UploadedFileRepository) Forget(ctx context.Context, fileID uuid.UUID) error {
	return r.remove(ctx, fileID, false)
}

func (r *UploadedFileRepository) remove(ctx context.Context, fileID uuid.UUID, cleanup bool) error {
	return r.Repository.Delete(ctx, func() error {
		if fileID == uuid.Nil() {
			return nil
		}

		if cleanup {
			uploadedFile := r.uploadedFiles.Find(fileID)
			if uploadedFile != nil && uploadedFile.upload != nil {
				// Keep the entry in the cache if the file could not be removed:
				// the file is still in the temporary directory, so the entry must keep tracking it.
				// Cleanup logs the failure itself.
				if err := uploadedFile.upload.Cleanup(); err != nil {
					return err
				}
			}
		}

		r.uploadedFiles.Delete(fileID)
		return nil
	})
}

// FindByDownloadID retrieves a mediaDownload by its fileID from the repository.
func (r *UploadedFileRepository) Find(
	ctx context.Context, fileID uuid.UUID,
) (pservices.UploadedFile, memsimple.CacheStatus, error) {
	find := func() (*uploadedFile, memsimple.CacheStatus, error) {
		uploadedFile, status := r.uploadedFiles.FindWithStatus(fileID)
		return uploadedFile, status, nil
	}

	uploadedFile, status, err := r.Repository.FindWithStatus(ctx, find)

	if uploadedFile == nil {
		return nil, status, err
	}

	return uploadedFile.upload, status, err
}

// Checks if a mediaDownload exists by its fileID.
func (r *UploadedFileRepository) Exists(ctx context.Context, fileID uuid.UUID) (bool, error) {
	exists := func() (bool, error) {
		return r.uploadedFiles.Exists(fileID), nil
	}
	return r.Repository.Exists(ctx, exists)
}

// CleanExpired cleans expired entries from the repository.
func (r *UploadedFileRepository) CleanExpired(ctx context.Context) error {
	// Define a clean function to remove expired entries from the cache.
	clean := func() error {
		r.uploadedFiles.CleanExpiredWithDeleteHook(
			func(key uuid.UUID, value *uploadedFile) {
				if value != nil && value.upload != nil {
					value.upload.Cleanup()
				}
			},
		)
		return nil
	}
	// Call the base repository's CleanExpired method with the custom clean function.
	return r.Repository.CleanExpired(ctx, clean)
}

// Clear clears all entries from the cache.
func (r *UploadedFileRepository) Clear(ctx context.Context) error {
	clear := func() error {
		r.uploadedFiles.ClearWithDeleteHook(
			func(key uuid.UUID, value *uploadedFile) {
				if value != nil && value.upload != nil {
					value.upload.Cleanup()
				}
			},
		)
		return nil
	}

	return r.Repository.Clear(ctx, clear)
}

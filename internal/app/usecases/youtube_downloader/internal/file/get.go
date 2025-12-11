package fileuc

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (uc *File) FindByFileId(ctx context.Context, fileId uuid.UUID, checkNotFound bool) (*ddownload.File, error) {
	file, err := uc.FileRep.FindByFileId(ctx, fileId)
	if err != nil {
		uc.logger.Warn("Failed to find record", "error", err)
		return nil, err
	}

	if checkNotFound && file == nil {
		uc.logger.Warn("Record not found", "fileId", fileId)
		return nil, errors.New("record not found")
	}

	return file, err
}

func (uc *File) GetAll(ctx context.Context, includeDeleted bool) ([]*ddownload.File, error) {
	file, err := uc.FileRep.GetAll(ctx, includeDeleted)
	if err != nil {
		uc.logger.Warn("Failed to get files", "error", err)
		return nil, err
	}

	return file, err
}

func (uc *File) GetAllFullNames(ctx context.Context, includeDeleted bool) ([]string, error) {
	names, err := uc.FileRep.GetAllFullNames(ctx, includeDeleted)
	if err != nil {
		uc.logger.Warn("Failed to get fullNames", "error", err)
		return nil, err
	}

	return names, nil
}

func (uc *File) GetBeforeTime(ctx context.Context, before time.Time, limit uint64) ([]*ddownload.File, error) {
	file, err := uc.FileRep.GetBeforeTime(ctx, before, limit)
	if err != nil {
		uc.logger.Warn("Failed to get files", "error", err)
		return nil, err
	}

	return file, err
}

func (uc *File) GetByStatus(ctx context.Context, status dtypes.FileStatus) ([]*ddownload.File, error) {
	file, err := uc.FileRep.GetByStatus(ctx, status)
	if err != nil {
		uc.logger.Warn("Failed to get files", "error", err)
		return nil, err
	}

	return file, err
}

func (uc *File) GetByPartialHash(ctx context.Context, hash string) ([]*ddownload.File, error) {
	file, err := uc.FileRep.GetByPartialHash(ctx, hash)
	if err != nil {
		uc.logger.Warn("Failed to get files", "error", err)
		return nil, err
	}

	return file, err
}

func (uc *File) GetWithoutPartialHash(ctx context.Context) ([]*ddownload.File, error) {
	var files []*ddownload.File

	gFiles, err := uc.FileRep.GetWithoutPartialHash(ctx)
	if err != nil {
		uc.logger.Warn("Failed to get files", "error", err)
		return nil, err
	}

	if len(gFiles) > 0 {
		files = make([]*ddownload.File, 0, len(gFiles))
		for _, file := range gFiles {
			if file.FullName == "" {
				continue
			}
			files = append(files, file)
		}
	}

	return files[:len(files):len(files)], nil
}

func (uc *File) GetDuplicateHashes(ctx context.Context) ([]string, error) {
	hashes, err := uc.FileRep.GetDuplicateHashes(ctx)
	if err != nil {
		uc.logger.Warn("Failed to get dublicate hashes", "error", err)
		return nil, err
	}

	return hashes, nil
}

func (uc *File) GetDeleted(ctx context.Context, from, to *time.Time) ([]*ddownload.File, error) {
	files, err := uc.FileRep.GetDeleted(ctx, from, to)
	if err != nil {
		uc.logger.Warn("Failed to get deleted", "fromDate", from, "toDate", to, "error", err)
		return nil, err
	}

	return files, nil
}

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
		uc.logger.Debug("Error finding record", "error", err)
		return nil, err
	}

	if checkNotFound && file == nil {
		err := errors.New("record not found")
		uc.logger.Debug("Record not found", "fileId", fileId)
		return nil, err
	}

	return file, err
}

func (uc *File) GetAll(ctx context.Context) ([]*ddownload.File, error) {
	file, err := uc.FileRep.GetAll(ctx)
	if err != nil {
		uc.logger.Debug("Error get files", "error", err)
		return nil, err
	}

	return file, err
}

func (uc *File) GetBeforeTime(ctx context.Context, before time.Time, limit uint64) ([]*ddownload.File, error) {
	file, err := uc.FileRep.GetBeforeTime(ctx, before, limit)
	if err != nil {
		uc.logger.Debug("Error get files", "error", err)
		return nil, err
	}

	return file, err
}

func (uc *File) GetByStatus(ctx context.Context, status dtypes.FileStatus) ([]*ddownload.File, error) {
	file, err := uc.FileRep.GetByStatus(ctx, status)
	if err != nil {
		uc.logger.Debug("Failed to get files", "error", err)
		return nil, err
	}

	return file, err
}

func (uc *File) GetByPartialHash(ctx context.Context, hash string) ([]*ddownload.File, error) {
	file, err := uc.FileRep.GetByPartialHash(ctx, hash)
	if err != nil {
		uc.logger.Debug("Failed to get files", "error", err)
		return nil, err
	}

	return file, err
}

func (uc *File) GetWithoutPartialHash(ctx context.Context) ([]*ddownload.File, error) {
	var files []*ddownload.File

	gFiles, err := uc.FileRep.GetWithoutPartialHash(ctx)
	if err != nil {
		uc.logger.Debug("Failed to get files", "error", err)
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
		uc.logger.Debug("Failed to get dublicate hashes")
		return nil, err
	}

	return hashes, nil
}

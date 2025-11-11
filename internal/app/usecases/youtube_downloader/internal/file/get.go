package file

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (uc *File) FindFileById(ctx context.Context, fileId uuid.UUID, checkNotFound bool) (*ddownload.File, error) {
	file, err := uc.fileRep.FindByFileId(ctx, fileId)
	if err != nil {
		uc.logger.Error("Error finding record", "error", err)
		return nil, err
	}

	if checkNotFound && file == nil {
		uc.logger.Error("Record not found", "fileId", fileId)
		return nil, err
	}

	return file, err
}

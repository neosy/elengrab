package fileuc

import (
	"context"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (uc *File) Update(ctx context.Context, file *ddownload.File) error {
	err := uc.fileRep.Update(ctx, file)
	if err != nil {
		uc.logger.Warn("Update record error", "error", err)
		return err
	}

	uc.saveToDownloadStateCache(ctx, file.FileID)

	return err
}

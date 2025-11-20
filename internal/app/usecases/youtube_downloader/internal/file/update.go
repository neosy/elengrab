package fileuc

import (
	"context"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (uc *File) Update(ctx context.Context, file *ddownload.File) error {
	err := uc.FileRep.Update(ctx, file)
	if err != nil {
		uc.logger.Debug("Update record error", "error", err)
		return err
	}
	return err
}

package file

import (
	"context"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (uc *File) Update(ctx context.Context, file *ddownload.File) error {
	err := uc.fileRep.Update(ctx, file)
	if err != nil {
		uc.logger.Error("Update record error", "error", err)
		return err
	}
	return err
}

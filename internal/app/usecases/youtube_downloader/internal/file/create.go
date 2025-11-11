package file

import (
	"context"
	"errors"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (uc *File) Create(ctx context.Context, file *ddownload.File) error {
	if file == nil {
		uc.logger.Error("Nil pointer in function")
		return errors.New("function parameter is a null pointer")
	}

	file.Status = dtypes.FileStatusNew

	err := uc.fileRep.Insert(ctx, file)
	if err != nil {
		uc.logger.Error(
			"Failed to insert record into repository",
			"error", err,
		)
		return err
	}

	return nil
}

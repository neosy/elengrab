package file

import (
	"context"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
)

func (uc *File) Patch(ctx context.Context, fileId uuid.UUID, patch *dto.FileInfoPatch) error {
	file, err := uc.FindFileById(ctx, fileId, true)
	if err != nil {
		return err
	}

	dto.PatchToFileDomain(patch, file)

	err = uc.Update(ctx, file)
	if err != nil {
		return err
	}

	return nil
}

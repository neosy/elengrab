package filestatus

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// updateStatus marking the status
func (uc *FileStatus) updateStatus(
	ctx context.Context,
	fileId uuid.UUID,
	toStatus dtypes.FileStatus,
	updateFieldsFunc func(file *ddownload.File),
) error {
	file, err := uc.file.FindByFileId(ctx, fileId, true)
	if err != nil {
		return err
	}

	err = uc.statusSetter.SetStatus(file, toStatus)
	if err != nil {
		uc.logger.Warn(
			"Failed to update status",
			"fileId", fileId,
			"error", err,
		)
		return err
	}

	// Update fields
	if updateFieldsFunc != nil {
		updateFieldsFunc(file)
	}

	// Update in the repository
	err = uc.file.Update(ctx, file)
	if err != nil {
		uc.logger.Warn(
			"Failed to update file in the repository",
			"fileId", fileId,
			"error", err,
		)

		return err
	}

	return nil
}

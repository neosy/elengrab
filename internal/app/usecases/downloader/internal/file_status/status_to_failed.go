package filestatus

import (
	"context"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// Failed set status to failed
func (s *FileStatus) Failed(
	ctx context.Context,
	fileID uuid.UUID,
	patch *dto.FileInfoPatch,
	message *string,
) error {
	updateFieldsFunc := func(file *ddownload.File) {
		dto.PatchToFileDomain(patch, file)
		file.ErrorMessage = message
	}

	task, err := s.dlTask.FindByFileID(ctx, fileID, true)
	if err != nil {
		return err
	}

	err = s.dlTaskStatus.Failed(ctx, task.TaskID)
	if err != nil {
		return err
	}

	return s.updateStatus(
		ctx,
		fileID,
		dtypes.FileStatusFailed,
		updateFieldsFunc,
	)
}

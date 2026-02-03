package filestatus

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// Failed set status to done
func (s *FileStatus) New(
	ctx context.Context,
	fileID uuid.UUID,
) error {
	updateFieldsFunc := func(file *ddownload.File) {
		file.ErrorMessage = nil
	}

	task, err := s.dlTask.FindByFileId(ctx, fileID, true)
	if err != nil {
		return err
	}

	err = s.dlTaskStatus.New(ctx, task.TaskId)
	if err != nil {
		return err
	}

	return s.updateStatus(
		ctx,
		fileID,
		dtypes.FileStatusNew,
		updateFieldsFunc,
	)
}

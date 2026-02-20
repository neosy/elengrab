package filestatus

import (
	"context"

	"github.com/google/uuid"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// Pending set status to working
func (uc *FileStatus) Working(
	ctx context.Context,
	fileId uuid.UUID,
	taskId uuid.UUID,
	workerId uint64,
) error {
	err := uc.dlTaskStatus.Working(ctx, taskId, workerId)
	if err != nil {
		return err
	}

	return uc.updateStatus(
		ctx,
		fileId,
		dtypes.FileStatusWorking,
		nil,
	)
}

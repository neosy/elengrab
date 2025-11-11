package filestatus

import (
	"context"

	"github.com/google/uuid"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// Pending set status to working
func (s *FileStatus) Working(
	ctx context.Context,
	fileId uuid.UUID,
) error {
	return s.updateStatus(
		ctx,
		fileId,
		dtypes.FileStatusWorking,
		nil,
	)
}

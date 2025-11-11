package filestatus

import (
	"context"

	"github.com/google/uuid"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// Pending set status to pending
func (s *FileStatus) Pending(
	ctx context.Context,
	fileId uuid.UUID,
) error {
	return s.updateStatus(
		ctx,
		fileId,
		dtypes.FileStatusPending,
		nil,
	)
}

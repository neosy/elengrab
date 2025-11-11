package statussetter

import (
	"fmt"

	dfilestatus "github.com/neosy/elengrab/internal/domain/download/file_status"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// checkSetStatus validates whether the status transition from fromStatus to toStatus is allowed.
func (u *FileStatusSetter) checkSetStatus(
	fromStatus dtypes.FileStatus,
	toStatus dtypes.FileStatus,
) error {
	_, exists := dfilestatus.SelectAllowedStatusMap()[fromStatus][toStatus]
	if exists {
		return nil
	}

	u.logger.Error(
		"Invalid status transition",
		"fromStatus", fromStatus,
		"toStatus", toStatus,
	)

	return fmt.Errorf("invalid status transition: %v -> %v", fromStatus, toStatus)
}

package statussetter

import (
	"fmt"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
	dfilestatus "github.com/neosy/elengrab/internal/domain/types/download_task_status"
)

// checkSetStatus validates whether the status transition from fromStatus to toStatus is allowed.
func (u *DownloadTaskStatusSetter) checkSetStatus(
	fromStatus dtypes.DownloadTaskStatus,
	toStatus dtypes.DownloadTaskStatus,
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

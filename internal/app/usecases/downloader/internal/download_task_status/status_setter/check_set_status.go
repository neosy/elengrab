package statussetter

import (
	"fmt"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
	dmediataskstatus "github.com/neosy/elengrab/internal/domain/types/media_task_status"
)

// checkSetStatus validates whether the status transition from fromStatus to toStatus is allowed.
func (u *StatusSetter) checkSetStatus(
	fromStatus dtypes.MediaTaskStatus,
	toStatus dtypes.MediaTaskStatus,
) error {
	_, exists := dmediataskstatus.SelectAllowedStatusMap()[fromStatus][toStatus]
	if exists {
		return nil
	}

	u.logger.Warn(
		"Invalid status transition",
		"fromStatus", fromStatus,
		"toStatus", toStatus,
	)

	return fmt.Errorf("invalid status transition: %v -> %v", fromStatus, toStatus)
}

package statussetter

import (
	"fmt"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
	mdlstatus "github.com/neosy/elengrab/internal/domain/types/media_download_status"
)

// checkSetStatus validates whether the status transition from fromStatus to toStatus is allowed.
func (u *MediaDownloadStatusSetter) checkSetStatus(
	fromStatus dtypes.MediaDownloadStatus,
	toStatus dtypes.MediaDownloadStatus,
) error {
	_, exists := mdlstatus.SelectAllowedStatusMap()[fromStatus][toStatus]
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

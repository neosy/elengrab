package dfilestatus

import dtypes "github.com/neosy/elengrab/internal/domain/types"

// MapAllowedStatus defines allowed transitions between statuses:
// key — current status, value — set of allowed next statuses.
type AllowedStatusMap = map[dtypes.DownloadTaskStatus]map[dtypes.DownloadTaskStatus]struct{}

var (
	allowedStatusMap = AllowedStatusMap{
		dtypes.DownloadTaskStatusPending: {
			dtypes.DownloadTaskStatusWorking: {},
		},
	}
)

// SelectAllowedStatusMap returns the allowed status transitions map based.
func SelectAllowedStatusMap() AllowedStatusMap {
	return allowedStatusMap
}

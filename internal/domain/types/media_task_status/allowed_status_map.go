package dmediataskstatus

import dtypes "github.com/neosy/elengrab/internal/domain/types"

// MapAllowedStatus defines allowed transitions between statuses:
// key — current status, value — set of allowed next statuses.
type AllowedStatusMap = map[dtypes.MediaTaskStatus]map[dtypes.MediaTaskStatus]struct{}

var (
	allowedStatusMap = AllowedStatusMap{
		dtypes.MediaTaskStatusNew: {
			dtypes.MediaTaskStatusPending: {},
		},
		dtypes.MediaTaskStatusPending: {
			dtypes.MediaTaskStatusWorking: {},
			dtypes.MediaTaskStatusFailed:  {},
		},
		dtypes.MediaTaskStatusWorking: {
			dtypes.MediaTaskStatusFailed: {},
		},
		dtypes.MediaTaskStatusFailed: {
			dtypes.MediaTaskStatusNew: {},
		},
	}
)

// SelectAllowedStatusMap returns the allowed status transitions map based.
func SelectAllowedStatusMap() AllowedStatusMap {
	return allowedStatusMap
}

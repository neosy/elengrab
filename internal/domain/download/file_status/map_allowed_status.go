package dfilestatus

import dtypes "github.com/neosy/elengrab/internal/domain/types"

// MapAllowedStatus defines allowed transitions between order statuses:
// key — current status, value — set of allowed next statuses.
type AllowedStatusMap = map[dtypes.FileStatus]map[dtypes.FileStatus]struct{}

var (
	mapAllowedStatus = AllowedStatusMap{
		dtypes.FileStatusNew: {
			dtypes.FileStatusPending: {},
			dtypes.FileStatusFailed:  {},
		},
		dtypes.FileStatusPending: {
			dtypes.FileStatusWorking: {},
			dtypes.FileStatusFailed:  {},
		},
		dtypes.FileStatusWorking: {
			dtypes.FileStatusDone:   {},
			dtypes.FileStatusFailed: {},
		},
	}
)

// SelectAllowedStatusMap returns the allowed order status transitions map based.
func SelectAllowedStatusMap() AllowedStatusMap {
	return mapAllowedStatus
}

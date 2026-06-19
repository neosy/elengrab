package mdlstatus

import dtypes "github.com/neosy/elengrab/internal/domain/types"

// MapAllowedStatus defines allowed transitions between statuses:
// key — current status, value — set of allowed next statuses.
type AllowedStatusMap = map[dtypes.MediaDownloadStatus]map[dtypes.MediaDownloadStatus]struct{}

var (
	allowedStatusMap = AllowedStatusMap{
		dtypes.MediaDownloadStatusNew: {
			dtypes.MediaDownloadStatusPending: {},
			dtypes.MediaDownloadStatusFailed:  {},
		},
		dtypes.MediaDownloadStatusPending: {
			dtypes.MediaDownloadStatusWorking: {},
			dtypes.MediaDownloadStatusFailed:  {},
		},
		dtypes.MediaDownloadStatusWorking: {
			dtypes.MediaDownloadStatusDone:   {},
			dtypes.MediaDownloadStatusFailed: {},
		},
		dtypes.MediaDownloadStatusFailed: {
			dtypes.MediaDownloadStatusNew: {},
		},
		dtypes.MediaDownloadStatusDone: {
			dtypes.MediaDownloadStatusRefreshing: {},
		},
		dtypes.MediaDownloadStatusRefreshing: {
			dtypes.MediaDownloadStatusDone: {},
		},
	}
)

// SelectAllowedStatusMap returns the allowed status transitions map based.
func SelectAllowedStatusMap() AllowedStatusMap {
	return allowedStatusMap
}

package dto

type BroadcastEventType uint8

const (
	BroadcastEventTypeNone BroadcastEventType = iota
	BroadcastEventTypeDownloadAdd
	BroadcastEventTypeDownloadUpdate
	BroadcastEventTypeDownloadDelete
	BroadcastEventTypeProgressUpdate
	BroadcastEventTypeDownloadStartRefreshing
	BroadcastEventTypeSystemInfoUpdate
	BroadcastEventTypeNotification
)

var (
	// broadcastEventTypeMap implementation of a set for BroadcastEventType
	broadcastEventTypeMap = map[BroadcastEventType]string{
		BroadcastEventTypeNone:                    "none",
		BroadcastEventTypeDownloadAdd:             "download_add",
		BroadcastEventTypeDownloadUpdate:          "download_update",
		BroadcastEventTypeDownloadDelete:          "download_delete",
		BroadcastEventTypeProgressUpdate:          "progress_update",
		BroadcastEventTypeDownloadStartRefreshing: "start_refreshing",
		BroadcastEventTypeSystemInfoUpdate:        "system_info_update",
		BroadcastEventTypeNotification:            "notification",
	}

	broadcastSSEEventNameByType = map[BroadcastEventType]string{
		BroadcastEventTypeDownloadAdd:             "row-add",
		BroadcastEventTypeDownloadUpdate:          "row-update",
		BroadcastEventTypeDownloadDelete:          "row-delete",
		BroadcastEventTypeProgressUpdate:          "row-patch-field",
		BroadcastEventTypeDownloadStartRefreshing: "row-start-refreshing",
		BroadcastEventTypeSystemInfoUpdate:        "system-info-update",
		BroadcastEventTypeNotification:            "notification",
	}
)

// String returns the value as a string.
func (v BroadcastEventType) String() string {
	return broadcastEventTypeMap[v]
}

// SSEEventName returns the event name used for SSE messages.
func (v BroadcastEventType) SSEEventName() string {
	return broadcastSSEEventNameByType[v]
}

// Exists returns true if the BroadcastEventType is valid.
func (v BroadcastEventType) Exists() bool {
	_, exists := broadcastEventTypeMap[v]
	return exists
}

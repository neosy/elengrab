package dto

type BroadcastEventType uint8

const (
	BroadcastEventTypeNone BroadcastEventType = iota
	BroadcastEventTypeDownloadAdd
	BroadcastEventTypeDownloadUpdate
	BroadcastEventTypeProgressUpdate
	BroadcastEventTypeDownloadDelete
	BroadcastEventTypeSystemInfoUpdate
	BroadcastEventTypeNotification
)

var (
	// broadcastEventTypeMap implementation of a set for BroadcastEventType
	broadcastEventTypeMap = map[BroadcastEventType]string{
		BroadcastEventTypeNone:             "none",
		BroadcastEventTypeDownloadAdd:      "download_add",
		BroadcastEventTypeDownloadUpdate:   "download_update",
		BroadcastEventTypeDownloadDelete:   "download_delete",
		BroadcastEventTypeProgressUpdate:   "progress_update",
		BroadcastEventTypeSystemInfoUpdate: "system_info_update",
		BroadcastEventTypeNotification:     "notification",
	}
)

// String returns the value as a string.
func (v BroadcastEventType) String() string {
	return broadcastEventTypeMap[v]
}

// Exists returns true if the BroadcastEventType is valid.
func (v BroadcastEventType) Exists() bool {
	_, exists := broadcastEventTypeMap[v]
	return exists
}

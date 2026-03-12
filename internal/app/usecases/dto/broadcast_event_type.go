package dto

type BroadcastEventType uint8

const (
	BroadcastEventTypeNone BroadcastEventType = iota
	BroadcastEventTypeFileAdd
	BroadcastEventTypeFileUpdate
	BroadcastEventTypeProgressUpdate
	BroadcastEventTypeFileDelete
)

var (
	// broadcastEventTypeMap implementation of a set for BroadcastEventType
	broadcastEventTypeMap = map[BroadcastEventType]string{
		BroadcastEventTypeNone:           "none",
		BroadcastEventTypeFileAdd:        "file_add",
		BroadcastEventTypeFileUpdate:     "file_update",
		BroadcastEventTypeFileDelete:     "file_delete",
		BroadcastEventTypeProgressUpdate: "progress_update",
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

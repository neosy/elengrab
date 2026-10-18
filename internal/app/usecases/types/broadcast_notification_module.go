package types

type BroadcastNotificationModule uint8

const (
	BroadcastNotificationModuleNone BroadcastNotificationModule = iota
	BroadcastNotificationModuleGrabForm
	BroadcastNotificationModuleUploadFile
	BroadcastNotificationModuleResultRow
)

var (
	// broadcastNotificationModuleValues implementation of a set for BroadcastNotificationModule
	broadcastNotificationModuleValues = map[BroadcastNotificationModule]string{
		BroadcastNotificationModuleNone:       "none",
		BroadcastNotificationModuleGrabForm:   "grab-form",
		BroadcastNotificationModuleUploadFile: "upload-file",
		BroadcastNotificationModuleResultRow:  "result-row",
	}
)

// String returns the value as a string.
func (v BroadcastNotificationModule) String() string {
	return broadcastNotificationModuleValues[v]
}

// Exists returns true if the BroadcastNotificationModule is valid.
func (v BroadcastNotificationModule) Exists() bool {
	_, exists := broadcastNotificationModuleValues[v]
	return exists
}

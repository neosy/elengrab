package dto

type BroadcastNotificationModule uint8

const (
	BroadcastNotificationModuleNone BroadcastNotificationModule = iota
	BroadcastNotificationModuleGrabForm
	BroadcastNotificationModuleResultRow
)

var (
	// broadcastNotificationModuleMap implementation of a set for BroadcastNotificationModule
	broadcastNotificationModuleMap = map[BroadcastNotificationModule]string{
		BroadcastNotificationModuleNone:      "none",
		BroadcastNotificationModuleGrabForm:  "grab-form",
		BroadcastNotificationModuleResultRow: "result-row",
	}
)

// String returns the value as a string.
func (v BroadcastNotificationModule) String() string {
	return broadcastNotificationModuleMap[v]
}

// Exists returns true if the BroadcastNotificationModule is valid.
func (v BroadcastNotificationModule) Exists() bool {
	_, exists := broadcastNotificationModuleMap[v]
	return exists
}

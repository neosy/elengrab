package dto

type BroadcastNotificationType uint8

const (
	BroadcastNotificationTypeNone BroadcastNotificationType = iota
	BroadcastNotificationTypeError
	BroadcastNotificationTypeWarning
	BroadcastNotificationTypeInfo
)

var (
	// broadcastNotificationTypeMap implementation of a set for BroadcastNotificationType
	broadcastNotificationTypeMap = map[BroadcastNotificationType]string{
		BroadcastNotificationTypeNone:    "none",
		BroadcastNotificationTypeError:   "error",
		BroadcastNotificationTypeWarning: "warning",
		BroadcastNotificationTypeInfo:    "info",
	}
)

// String returns the value as a string.
func (v BroadcastNotificationType) String() string {
	return broadcastNotificationTypeMap[v]
}

// Exists returns true if the BroadcastNotificationType is valid.
func (v BroadcastNotificationType) Exists() bool {
	_, exists := broadcastNotificationTypeMap[v]
	return exists
}

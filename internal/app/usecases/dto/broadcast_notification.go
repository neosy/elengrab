package dto

type BroadcastNotification struct {
	Module  BroadcastNotificationModule
	Type    BroadcastNotificationType
	Message string
}

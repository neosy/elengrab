package dto

import "github.com/neosy/elengrab/internal/app/usecases/types"

type BroadcastNotification struct {
	Module  types.BroadcastNotificationModule
	Type    types.BroadcastNotificationType
	Message string
}

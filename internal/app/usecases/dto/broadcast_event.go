package dto

import "github.com/neosy/elengrab/internal/app/usecases/types"

type BroadcastEvent struct {
	Type types.BroadcastEventType
	Data any
}

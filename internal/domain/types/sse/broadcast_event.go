package dsse

import (
	"time"

	"github.com/neosy/elengrab/internal/app/usecases/types"
)

type BroadcastEvent struct {
	Type      types.BroadcastEventType
	Data      any
	Timestamp int64
}

// NewBroadcastEvent creates a new BroadcastEvent
func NewBroadcastEvent(eventType types.BroadcastEventType, data any) BroadcastEvent {
	return BroadcastEvent{
		Type:      eventType,
		Data:      data,
		Timestamp: time.Now().UnixMicro(),
	}
}

// EmptyBroadcastEvent creates a new BroadcastEvent with no type or data
func EmptyBroadcastEvent() BroadcastEvent {
	return BroadcastEvent{}
}

// IsEmpty returns true if the event is empty (no type or data)
func (e BroadcastEvent) IsEmpty() bool {
	return e.Type == types.BroadcastEventTypeNone || e.Data == nil
}

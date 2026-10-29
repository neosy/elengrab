package dsse

import (
	"fmt"
	"uuid"

	"github.com/neosy/elengrab/internal/pkg/idcodec"
)

// EventID identifies an event inside a replay log.
type EventID struct {
	LogID  uuid.UUID
	LogNum int
}

// EventIDNil returns a nil SSEEventID.
func EventIDNil() EventID {
	return EventID{}
}

// NewEventID creates a new SSEEventID with the given logId and logNum.
func NewEventID(logId uuid.UUID, logNum int) EventID {
	return EventID{
		LogID:  logId,
		LogNum: logNum,
	}
}

// String returns the wire form used in the SSE "id:" field.
func (id *EventID) String() string {
	if id == nil {
		return ""
	}

	return fmt.Sprintf("%s:%d", idcodec.EncodeUUIDBase64URL(id.LogID), id.LogNum)
}

// IsZero returns true if the SSEEventID is nil or has a zero LogID.
func (id *EventID) IsZero() bool {
	return id == nil || id.LogID == uuid.Nil()
}

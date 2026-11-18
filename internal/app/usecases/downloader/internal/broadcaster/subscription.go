package broadcaster

import (
	"uuid"

	dsse "github.com/neosy/elengrab/internal/domain/types/sse"
)

type subscription struct {
	connectionID uuid.UUID
	userID       uuid.UUID
	eventCh      chan dsse.BroadcastEvent
}

func (s subscription) EventCh() chan dsse.BroadcastEvent {
	return s.eventCh
}

func (c subscription) ConnectionID() uuid.UUID {
	return c.connectionID
}

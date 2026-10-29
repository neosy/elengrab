package broadcaster

import (
	"uuid"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
	dsse "github.com/neosy/elengrab/internal/domain/types/sse"
)

type subscription struct {
	connectionID uuid.UUID
	roles        dtypes.UserRoleIDs
	eventCh      chan dsse.BroadcastEvent
}

func (s subscription) EventCh() chan dsse.BroadcastEvent {
	return s.eventCh
}

func (c subscription) ConnectionID() uuid.UUID {
	return c.connectionID
}

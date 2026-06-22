package broadcaster

import (
	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type subscription struct {
	connectionID uuid.UUID
	roles        dtypes.UserRoleIDs
	eventCh      chan dto.BroadcastEvent
}

func (s subscription) EventCh() chan dto.BroadcastEvent {
	return s.eventCh
}

func (c subscription) ConnectionID() uuid.UUID {
	return c.connectionID
}

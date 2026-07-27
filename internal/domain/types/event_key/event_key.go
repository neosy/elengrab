package eventkey

import (
	"fmt"

	"github.com/google/uuid"
)

type EventKey struct {
	Type EventKeyType
	ID   string
}

func NewEventKey(t EventKeyType, id string) EventKey {
	return EventKey{
		Type: t,
		ID:   id,
	}
}

func NewEventKeyUserID(userID uuid.UUID) EventKey {
	return NewEventKey(EventKeyTypeUser, userID.String())
}

func NewEventKeySessionID(sessionID uuid.UUID) EventKey {
	return NewEventKey(EventKeyTypeSession, sessionID.String())
}

func (e EventKey) UUID() uuid.UUID {
	data := fmt.Sprintf("%s:%s", e.Type.String(), e.ID)
	return uuid.NewSHA1(uuid.NameSpaceOID, []byte(data))
}

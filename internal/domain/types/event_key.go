package dtypes

import (
	"fmt"

	"github.com/google/uuid"
)

type EventKeyType uint8

const (
	EventKeyTypeUser EventKeyType = iota
	EventKeyTypeSession
)

var (
	enentKeyNameByType = map[EventKeyType]string{
		EventKeyTypeUser:    "user",
		EventKeyTypeSession: "session",
	}
)

func (v EventKeyType) String() string {
	return enentKeyNameByType[v]
}

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

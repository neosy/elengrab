package dtypes

import "github.com/google/uuid"

type EventKeyType uint8

const (
	EventKeyTypeUser EventKeyType = iota
	EventKeyTypeAnonSession
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

func NewEventKeyUserID(id uuid.UUID) EventKey {
	return NewEventKeyUser(id.String())
}

func NewEventKeyUser(id string) EventKey {
	return EventKey{
		Type: EventKeyTypeUser,
		ID:   id,
	}
}

func NewEventKeyAnonSessionID(id uuid.UUID) EventKey {
	return NewEventKeyAnonSession(id.String())
}

func NewEventKeyAnonSession(id string) EventKey {
	return EventKey{
		Type: EventKeyTypeAnonSession,
		ID:   id,
	}
}

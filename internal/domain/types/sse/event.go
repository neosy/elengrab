package dsse

import "uuid"

type Event struct {
	Name string
	Data []byte

	Num   int
	Index string // dedup key: "userKey:dataHash"
}

// NewEvent creates a new SSE event with the given name and data.
func NewEvent(name string, data []byte) Event {
	return Event{
		Name: name,
		Data: data,
	}
}

// EventID returns the EventID for the given logID and the event's number.
func (e *Event) EventID(logID uuid.UUID) EventID {
	if e == nil {
		return EventIDNil()
	}

	return NewEventID(logID, e.Num)
}

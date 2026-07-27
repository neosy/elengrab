package eventkey

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

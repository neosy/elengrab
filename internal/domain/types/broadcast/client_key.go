package broadcast

import (
	"github.com/google/uuid"
	eventkey "github.com/neosy/elengrab/internal/domain/types/event_key"
)

type ClientKey struct {
	ConnectionID uuid.UUID
	EventKey     eventkey.EventKey
}

func BuildClientKey(eventKey eventkey.EventKey) ClientKey {
	return ClientKey{
		ConnectionID: uuid.New(),
		EventKey:     eventKey,
	}
}

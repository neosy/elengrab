package persistence

import (
	"context"
	"uuid"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/cache/memory"
)

type SSEEventLogCacheRepository interface {
	memory.CacheRepository

	// Push adds a new SSE event to the log for a given key in the repository.
	Push(
		ctx context.Context,
		key string,
		event *dtypes.SSEEvent,
	) (logID uuid.UUID, logNum int, err error)

	// GetEvents retrieves all events for a given key from the repository.
	//
	// The returned slice is a copy, but the events it points to are shared
	// and MUST NOT be modified by the caller.
	GetEvents(ctx context.Context, key string) ([]*dtypes.SSEEvent, error)

	// GetEventsAfter retrieves events with numbers greater than lastNum from the event log.
	//
	// The returned slice is a copy, but the events it points to are shared
	// and MUST NOT be modified by the caller.
	GetEventsAfter(
		ctx context.Context,
		key string,
		logID uuid.UUID,
		lastNum int,
	) (events []*dtypes.SSEEvent, resync bool, err error)

	// GetPosition returns the logID and the last event number of the log for a key.
	GetPosition(
		ctx context.Context,
		key string,
	) (logID uuid.UUID, lastNum int)

	// CleanExpired cleans expired entries from the repository.
	CleanExpired(ctx context.Context) error
}

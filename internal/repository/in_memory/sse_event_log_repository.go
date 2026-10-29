package inmemory

import (
	"context"
	"slices"
	"time"
	"uuid"

	dsse "github.com/neosy/elengrab/internal/domain/types/sse"
	memsimple "github.com/neosy/elengrab/internal/pkg/cache/memory/simple"
)

const defaultLenSseEventLog = 1000

type sseEventLog struct {
	id      uuid.UUID // changes on every recreation of the log
	lastNum int
	items   []*dsse.Event

	// uniqueIndex tracks uniqueness of event data: "userKey:dataHash"
	uniqueIndex map[string]int
}

func newSseEventLog() *sseEventLog {
	return &sseEventLog{
		id:          uuid.New(),
		items:       make([]*dsse.Event, 0, defaultLenSseEventLog),
		uniqueIndex: make(map[string]int),
	}
}

// LastEventID returns the logID and the last event number of the log.
func (l *sseEventLog) LastEventID() dsse.EventID {
	if l == nil {
		return dsse.EventID{}
	}

	return dsse.EventID{
		LogID:  l.id,
		LogNum: l.lastNum,
	}
}

func (l *sseEventLog) Clone() *sseEventLog {
	newLog := *l
	// Events are immutable after publication, so sharing pointers is safe
	newLog.items = slices.Clone(l.items)
	return &newLog
}

type sseEventLogRepository struct {
	// Embeds the base Repository
	memsimple.Repository[sseEventLog]

	// Cache for storing SSE events keyed by event key
	cacheByKey memsimple.Cache[string, sseEventLog]
}

// newSseEventLogRepository creates a new instance of sseEventLogRepository
func newSseEventLogRepository(ttl time.Duration) *sseEventLogRepository {
	r := &sseEventLogRepository{
		cacheByKey: memsimple.NewCacheWithDeaultCloner[string, sseEventLog, *sseEventLog](),
	}
	r.Repository.Init(ttl)
	return r
}

func (r *sseEventLogRepository) Name() string {
	return "sseEventLog"
}

// Push adds a new SSE event to the log for a given key in the repository.
func (r *sseEventLogRepository) Push(
	ctx context.Context,
	key string,
	event *dsse.Event,
) (dsse.EventID, error) {
	var (
		eventIDNil dsse.EventID
		eventID    dsse.EventID
	)

	if key == "" || event == nil {
		return eventIDNil, nil
	}

	// Hash is computed once, outside the lock
	index := key + ":" + dataHashString(event.Data)

	save := func() error {
		var log *sseEventLog

		// The stored log is never modified, a new one is built instead
		old := r.cacheByKey.FindNoClone(key)

		if old != nil {
			lastNum, exists := old.uniqueIndex[index]
			if exists {
				eventID.LogID = old.id
				eventID.LogNum = lastNum
				return nil
			}

			log = new(*old)

			// Keep only the newest events
			start := max(0, len(old.items)-defaultLenSseEventLog+1)
			log.items = make([]*dsse.Event, 0, defaultLenSseEventLog)
			log.items = append(log.items, old.items[start:]...)

			// uniqueIndex is shared with the old log: it is only accessed under the mutex.
			// Evicted events must leave the index too
			for _, item := range old.items[:start] {
				delete(log.uniqueIndex, item.Index)
			}
		}

		if log == nil {
			log = newSseEventLog()
		}

		log.lastNum++

		event.Num = log.lastNum
		event.Index = index

		log.items = append(log.items, event)
		log.uniqueIndex[index] = log.lastNum

		r.cacheByKey.SaveNoClone(
			key,
			log,
			r.TTL(),
		)

		eventID.LogID = log.id
		eventID.LogNum = log.lastNum

		return nil
	}

	err := r.Repository.Save(ctx, save)
	if err != nil {
		return eventIDNil, err
	}

	return eventID, nil
}

// Delete removes the SSE event log for a given key from the repository.
func (r *sseEventLogRepository) Delete(ctx context.Context, key string) error {
	delete := func() error {
		if key != "" {
			r.cacheByKey.Delete(key)
		}
		return nil
	}
	return r.Repository.Delete(ctx, delete)
}

// GetEvents retrieves all events for a given key from the repository.
//
// The returned slice is a copy, but the events it points to are shared
// and MUST NOT be modified by the caller.
func (r *sseEventLogRepository) GetEvents(
	ctx context.Context,
	key string,
) ([]*dsse.Event, error) {
	log, err := r.Repository.Find(ctx, func() (*sseEventLog, error) {
		data := r.cacheByKey.Find(key)
		return data, nil
	})
	if err != nil {
		return nil, err
	}

	if log == nil {
		return nil, nil
	}

	return log.items, nil
}

// GetEventsAfter retrieves events with numbers greater than lastNum from the event log.
//
// The returned slice is a copy, but the events it points to are shared
// and MUST NOT be modified by the caller.
func (r *sseEventLogRepository) GetEventsAfter(
	ctx context.Context,
	key string,
	eventID dsse.EventID,
) (events []*dsse.Event, resync bool, err error) {
	log, err := r.Repository.Find(ctx, func() (*sseEventLog, error) {
		data := r.cacheByKey.Find(key)
		return data, nil
	})
	if err != nil {
		return nil, false, err
	}

	if log == nil {
		return nil, !eventID.IsZero(), nil
	}

	// Log was recreated: the client's numbers belong to another log
	if eventID.LogID != log.id {
		return nil, true, nil
	}

	lastNum := eventID.LogNum

	if lastNum > log.lastNum {
		return nil, true, nil
	}

	eventCount := log.lastNum - lastNum
	if eventCount <= 0 {
		return nil, false, nil
	}

	if eventCount > len(log.items) {
		return log.items, true, nil
	}

	startPos := len(log.items) - eventCount

	return log.items[startPos:], false, nil
}

// GetPosition returns the logID and the last event number of the log for a key.
func (r *sseEventLogRepository) GetPosition(
	ctx context.Context,
	key string,
) dsse.EventID {
	log, _ := r.Repository.Find(ctx, func() (*sseEventLog, error) {
		log := r.cacheByKey.FindNoClone(key)
		if log == nil {
			return nil, nil
		}

		return &sseEventLog{
			id:      log.id,
			lastNum: log.lastNum,
		}, nil
	})

	return log.LastEventID()
}

// CleanExpired cleans expired entries from the repository.
func (r *sseEventLogRepository) CleanExpired(ctx context.Context) error {
	// Define a clean function to remove expired entries from the cache.
	clean := func() error {
		r.cacheByKey.CleanExpired()
		return nil
	}
	// Call the base repository's CleanExpired method with the custom clean function.
	return r.Repository.CleanExpired(ctx, clean)
}

package inmemory

import (
	"context"
	"slices"
	"time"
	"uuid"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
	memsimple "github.com/neosy/elengrab/internal/pkg/cache/memory/simple"
)

const defaultLenSseEventLog = 10

type sseEventLog struct {
	id      uuid.UUID // changes on every recreation of the log
	lastNum int
	items   []*dtypes.SSEEvent

	// uniqueIndex tracks uniqueness of event data: "userKey:dataHash"
	uniqueIndex map[string]int
}

func newSseEventLog() *sseEventLog {
	return &sseEventLog{
		id:          uuid.New(),
		items:       make([]*dtypes.SSEEvent, 0, defaultLenSseEventLog),
		uniqueIndex: make(map[string]int),
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
	event *dtypes.SSEEvent,
) (logID uuid.UUID, logNum int, err error) {
	if key == "" || event == nil {
		return uuid.Nil(), 0, nil
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
				logID = old.id
				logNum = lastNum
				return nil
			}

			log = new(*old)

			// Keep only the newest events
			start := max(0, len(old.items)-defaultLenSseEventLog+1)
			log.items = make([]*dtypes.SSEEvent, 0, defaultLenSseEventLog)
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

		logID = log.id
		logNum = log.lastNum

		return nil
	}

	err = r.Repository.Save(ctx, save)
	if err != nil {
		return uuid.Nil(), 0, err
	}

	return logID, logNum, nil
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
) ([]*dtypes.SSEEvent, error) {
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
	logID uuid.UUID,
	lastNum int,
) (events []*dtypes.SSEEvent, resync bool, err error) {
	log, err := r.Repository.Find(ctx, func() (*sseEventLog, error) {
		data := r.cacheByKey.Find(key)
		return data, nil
	})
	if err != nil {
		return nil, false, err
	}

	if log == nil {
		return nil, logID != uuid.Nil(), nil
	}

	// Log was recreated: the client's numbers belong to another log
	if logID != log.id {
		return nil, true, nil
	}

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
) (logID uuid.UUID, lastNum int) {
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

	if log == nil {
		return uuid.Nil(), 0
	}

	return log.id, log.lastNum
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

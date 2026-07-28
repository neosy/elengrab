package mediawatch

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type statsUpdateKey struct {
	downloadID uuid.UUID
	userID     uuid.UUID
}

type statsUpdateQueue struct {
	mu    sync.RWMutex
	items map[statsUpdateKey]time.Duration
}

func newStatsUpdateQueue() statsUpdateQueue {
	return statsUpdateQueue{
		items: make(map[statsUpdateKey]time.Duration),
	}
}

func (q *statsUpdateQueue) add(downloadID uuid.UUID, userID uuid.UUID, mediaDuration time.Duration) {
	q.mu.Lock()
	q.items[statsUpdateKey{downloadID, userID}] = mediaDuration
	q.mu.Unlock()
}

func (q *statsUpdateQueue) drain() map[statsUpdateKey]time.Duration {
	q.mu.Lock()
	defer q.mu.Unlock()

	items := q.items
	q.items = make(map[statsUpdateKey]time.Duration)

	return items
}

func (q *statsUpdateQueue) count() int {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return len(q.items)
}

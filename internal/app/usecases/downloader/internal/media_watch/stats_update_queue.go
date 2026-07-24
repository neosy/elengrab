package mediawatch

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type statsUpdateQueue struct {
	mu    sync.RWMutex
	items map[uuid.UUID]time.Duration
}

func newStatsUpdateQueue() statsUpdateQueue {
	return statsUpdateQueue{
		items: make(map[uuid.UUID]time.Duration),
	}
}

func (q *statsUpdateQueue) add(downloadID uuid.UUID, mediaDuration time.Duration) {
	q.mu.Lock()
	q.items[downloadID] = mediaDuration
	q.mu.Unlock()
}

func (q *statsUpdateQueue) drain() map[uuid.UUID]time.Duration {
	q.mu.Lock()
	defer q.mu.Unlock()

	items := q.items
	q.items = make(map[uuid.UUID]time.Duration)

	return items
}

func (q *statsUpdateQueue) count() int {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return len(q.items)
}

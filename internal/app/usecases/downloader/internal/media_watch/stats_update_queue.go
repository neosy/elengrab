package mediawatch

import (
	"sync"
	"time"

	"uuid"

	"github.com/neosy/elengrab/internal/app/usecases/dto"
)

type statsUpdateKey struct {
	dto.MediaWatchUpdateKey
}

type downloadUserKey struct {
	downloadID uuid.UUID
	userID     uuid.UUID
}

type statsUpdateQueue struct {
	mu    sync.RWMutex
	items map[statsUpdateKey]time.Duration
}

func (k statsUpdateKey) downloadUserKey() downloadUserKey {
	return downloadUserKey{
		downloadID: k.DownloadID,
		userID:     k.UserID,
	}
}

func newStatsUpdateQueue() statsUpdateQueue {
	return statsUpdateQueue{
		items: make(map[statsUpdateKey]time.Duration),
	}
}

func (q *statsUpdateQueue) add(downloadID uuid.UUID, userID uuid.UUID, sessionID uuid.UUID, mediaDuration time.Duration) {
	q.mu.Lock()
	q.items[statsUpdateKey{
		DownloadID: downloadID,
		UserID:     userID,
		SessionID:  sessionID,
	}] = mediaDuration
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

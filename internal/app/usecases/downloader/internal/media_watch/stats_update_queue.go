package mediawatch

import (
	"sync"
	"time"

	"github.com/google/uuid"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
)

type statsUpdateKey struct {
	downloadID uuid.UUID
	userID     uuid.UUID
	sessionID  uuid.UUID
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
		downloadID: k.downloadID,
		userID:     k.userID,
	}
}

func (k statsUpdateKey) authCtx() dauth.AuthContext {
	var sessionID uuid.UUID
	if k.userID == uuid.Nil {
		sessionID = k.sessionID
	}

	return dauth.AuthContext{
		UserID:        k.userID,
		AnonSessionID: sessionID,
	}
}

func newStatsUpdateQueue() statsUpdateQueue {
	return statsUpdateQueue{
		items: make(map[statsUpdateKey]time.Duration),
	}
}

func (q *statsUpdateQueue) add(downloadID uuid.UUID, userID uuid.UUID, sessionID uuid.UUID, mediaDuration time.Duration) {
	q.mu.Lock()
	q.items[statsUpdateKey{downloadID, userID, sessionID}] = mediaDuration
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

package lock

import "sync"

type SQLiteLock struct {
	mu sync.RWMutex
}

func NewSQLiteLock() *SQLiteLock {
	return &SQLiteLock{}
}

func (l *SQLiteLock) Lock() {
	l.mu.Lock()
}

func (l *SQLiteLock) Unlock() {
	l.mu.Unlock()
}

func (l *SQLiteLock) RLock() {
	l.mu.RLock()
}

func (l *SQLiteLock) RUnlock() {
	l.mu.RUnlock()
}

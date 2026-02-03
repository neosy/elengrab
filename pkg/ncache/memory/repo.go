package nmemory

import (
	"sync"
	"time"
)

// Repository provides a thread-safe wrapper for operations on cached items.
// It handles locking and TTL for safe concurrent access.
type Repository[T any] struct {
	mu  *sync.RWMutex
	ttl time.Duration
}

// Init sets the default TTL for the repository items.
func (r *Repository[T]) Init(ttl time.Duration) {
	r.mu = &sync.RWMutex{}
	r.ttl = ttl
}

// TTL returns the default time-to-live duration for items in the repository.
func (r *Repository[T]) TTL() time.Duration {
	return r.ttl
}

// Save executes a write operation safely under a write lock.
func (r *Repository[T]) Save(fnSave func() error) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	return fnSave()
}

// Delete executes a deletion operation safely under a write lock.
func (r *Repository[T]) Delete(fnDelete func() error) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	return fnDelete()
}

// Find executes a read operation safely under a read lock.
func (r *Repository[T]) Find(fnFind func() (*T, error)) (*T, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return fnFind()
}

// Exists executes a read operation to check existence safely under a read lock.
func (r *Repository[T]) Exists(fnExists func() (bool, error)) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return fnExists()
}

// CleanExpired executes a cleanup operation safely under a write lock.
func (r *Repository[T]) CleanExpired(fnClean func() error) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	return fnClean()
}

// CopyAdapter converts a copy function into a CacheCopier that ignores its input.
func (r *Repository[T]) CopyAdapter(makeCopy func() *T) CacheCopier[T] {
	return func(*T) *T { return makeCopy() }
}

package nmemory

import (
	"sync"
	"time"
)

// Repository provides a thread-safe wrapper for operations on cached items.
// It handles locking and TTL for safe concurrent access.
type Repository[T any] struct {
	mu  sync.RWMutex
	ttl time.Duration
}

// Init sets the default TTL for the repository items.
func (r *Repository[T]) Init(ttl time.Duration) {
	r.ttl = ttl
}

// TTL returns the default time-to-live duration for items in the repository.
func (r *Repository[T]) TTL() time.Duration {
	return r.ttl
}

// Save executes a write operation safely under a write lock.
func (r *Repository[T]) Save(save func() error) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	return save()
}

// Delete executes a deletion operation safely under a write lock.
func (r *Repository[T]) Delete(delete func() error) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	return delete()
}

// Find executes a read operation safely under a read lock.
func (r *Repository[T]) Find(find func() (*T, error)) (*T, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return find()
}

// Exists executes a read operation to check existence safely under a read lock.
func (r *Repository[T]) Exists(exists func() (bool, error)) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return exists()
}

// CleanExpired executes a cleanup operation safely under a write lock.
func (r *Repository[T]) CleanExpired(clean func() error) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	return clean()
}

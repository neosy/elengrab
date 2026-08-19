package memsimple

import (
	"context"
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
// Use this for standard cache write operations that do not require a precomputed time value.
func (r *Repository[T]) Save(ctx context.Context, fnSave func() error) error {
	if !isTransactionContext(ctx) {
		r.mu.Lock()
		defer r.mu.Unlock()
	}

	return fnSave()
}

// SaveWithNow executes a write operation safely under a write lock.
// It is intended for high-performance write operations that use a precomputed
// current time value to avoid repeated time.Now() calls in hot paths.
func (r *Repository[T]) SaveWithNow(ctx context.Context, fnSave func() error) error {
	if !isTransactionContext(ctx) {
		r.mu.Lock()
		defer r.mu.Unlock()
	}

	return fnSave()
}

// Delete executes a deletion operation safely under a write lock.
func (r *Repository[T]) Delete(ctx context.Context, fnDelete func() error) error {
	if !isTransactionContext(ctx) {
		r.mu.Lock()
		defer r.mu.Unlock()
	}

	return fnDelete()
}

// Find executes a read operation safely under a read lock.
func (r *Repository[T]) Find(ctx context.Context, fnFind func() (*T, error)) (*T, error) {
	if !isTransactionContext(ctx) {
		r.mu.RLock()
		defer r.mu.RUnlock()
	}

	return fnFind()
}

// FindWithNow executes a read operation safely under a read lock.
// It is intended for high-performance cache lookups that provide a precomputed
// current time value to avoid repeated time.Now() calls in hot paths.
func (r *Repository[T]) FindWithNow(ctx context.Context, fnFind func() (*T, error)) (*T, error) {
	if !isTransactionContext(ctx) {
		r.mu.RLock()
		defer r.mu.RUnlock()
	}

	return fnFind()
}

// FindWithStatus executes a read operation safely under a read lock, returning the cache status.
func (r *Repository[T]) FindWithStatus(ctx context.Context, fnFind func() (*T, CacheStatus, error)) (*T, CacheStatus, error) {
	if !isTransactionContext(ctx) {
		r.mu.RLock()
		defer r.mu.RUnlock()
	}

	return fnFind()
}

// FindWithStatusNow executes a read operation safely under a read lock.
// It is intended for high-performance cache lookups that provide a precomputed
// current time value to avoid repeated time.Now() calls in hot paths.
func (r *Repository[T]) FindWithStatusNow(ctx context.Context, fnFind func() (*T, CacheStatus, error)) (*T, CacheStatus, error) {
	if !isTransactionContext(ctx) {
		r.mu.RLock()
		defer r.mu.RUnlock()
	}

	return fnFind()
}

// Exists executes a read operation to check existence safely under a read lock.
func (r *Repository[T]) Exists(ctx context.Context, fnExists func() (bool, error)) (bool, error) {
	if !isTransactionContext(ctx) {
		r.mu.RLock()
		defer r.mu.RUnlock()
	}

	return fnExists()
}

// ExistsWithNow executes a read operation to check existence safely under a read lock.
// It is intended for high-performance cache lookups that provide a precomputed
// current time value to avoid repeated time.Now() calls in hot paths.
func (r *Repository[T]) ExistsWithNow(ctx context.Context, fnExists func() (bool, error)) (bool, error) {
	if !isTransactionContext(ctx) {
		r.mu.RLock()
		defer r.mu.RUnlock()
	}

	return fnExists()
}

// ExistsWithStatus executes a read operation to check existence safely under a read lock, returning the cache status.
func (r *Repository[T]) ExistsWithStatus(ctx context.Context, fnExists func() (bool, CacheStatus, error)) (bool, CacheStatus, error) {
	if !isTransactionContext(ctx) {
		r.mu.RLock()
		defer r.mu.RUnlock()
	}

	return fnExists()
}

// ExistsWithStatusNow executes a read operation to check existence safely under a read lock.
// It is intended for high-performance cache lookups that provide a precomputed
// current time value to avoid repeated time.Now() calls in hot paths.
func (r *Repository[T]) ExistsWithStatusNow(ctx context.Context, fnExists func() (bool, CacheStatus, error)) (bool, CacheStatus, error) {
	if !isTransactionContext(ctx) {
		r.mu.RLock()
		defer r.mu.RUnlock()
	}

	return fnExists()
}

// CleanExpired executes a cleanup operation safely under a write lock.
func (r *Repository[T]) CleanExpired(ctx context.Context, fnClean func() error) error {
	if !isTransactionContext(ctx) {
		r.mu.Lock()
		defer r.mu.Unlock()
	}

	return fnClean()
}

// CopyAdapter converts a copy function into a CacheCopier that ignores its input.
func (r *Repository[T]) CopyAdapter(makeCopy func() *T) CacheCopier[T] {
	return func(*T) *T { return makeCopy() }
}

// Transaction executes fn under a repository-wide write lock.
func (r *Repository[T]) Transaction(fn func(context.Context) error) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	ctx := withTransactionContext(context.Background())

	return fn(ctx)
}

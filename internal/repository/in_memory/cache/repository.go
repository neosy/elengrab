package cache

import (
	"sync"
	"time"
)

type Repository[T any] struct {
	mu  sync.RWMutex
	ttl time.Duration
}

func (r *Repository[T]) Init(ttl time.Duration) {
	r.ttl = ttl
}

func (r *Repository[T]) TTL() time.Duration {
	return r.ttl
}

func (r *Repository[T]) Save(save func() error) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	return save()
}

func (r *Repository[T]) Delete(delete func() error) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	return delete()
}

func (r *Repository[T]) Find(find func() (*T, error)) (*T, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return find()
}

func (r *Repository[T]) Exists(exists func() (bool, error)) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return exists()
}

func (r *Repository[T]) CleanExpired(clean func() error) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	return clean()
}

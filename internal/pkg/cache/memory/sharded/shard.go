package memsharded

import "sync"

// shard represents a single partition of the cache map.
// Each shard holds a subset of keys and protects them with its own lock.
type shard[K comparable, T any] struct {
	mu    sync.RWMutex
	cache map[K]*Item[T]
}

// newShard creates an initialized shard instance.
func newShard[K comparable, T any]() *shard[K, T] {
	return &shard[K, T]{
		cache: make(map[K]*Item[T]),
	}
}

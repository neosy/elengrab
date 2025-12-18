package nmemory

import (
	"time"
)

// Cache represents a generic in-memory cache mapping keys to Items.
type Cache[K comparable, T any] map[K]*Item[T]

// Save stores a value in the cache with an optional copy function and TTL.
// If ttl > 0, the item expires after the specified duration.
func (c Cache[K, T]) Save(key K, value *T, fnCopy func(value *T) *T, ttl time.Duration) {
	if value == nil {
		return
	}

	valueCopy := value
	if fnCopy != nil {
		valueCopy = fnCopy(value)
	}

	expiresAt := time.Time{}
	if ttl > 0 {
		expiresAt = time.Now().Add(ttl)
	}

	c[key] = &Item[T]{
		value:     valueCopy,
		expiresAt: expiresAt,
	}
}

// Delete removes the item with the given key from the cache
func (c Cache[K, T]) Delete(key K) {
	if _, exists := c[key]; !exists {
		return
	}
	delete(c, key)
}

// Find returns the cached value for the given key, or nil if not found or expired.
// An optional copy function can be used to return a copy of the value.
func (c Cache[K, T]) Find(key K, fnCopy func(value *T) *T) *T {
	cacheData, exists := c[key]
	if !exists {
		return nil
	}

	if cacheData.Expired() {
		delete(c, key)
		return nil
	}

	valueCopy := cacheData.value
	if fnCopy != nil {
		valueCopy = fnCopy(cacheData.value)
	}

	return valueCopy
}

// Exists checks if the cache contains a non-expired item for the given key.
func (c Cache[K, T]) Exists(key K) bool {
	cacheData, exists := c[key]
	if !exists {
		return false
	}

	if cacheData.Expired() {
		delete(c, key)
		return false
	}

	return true
}

// CleanExpired removes all expired items from the cache.
func (c Cache[K, T]) CleanExpired() {
	now := time.Now()
	for key, cacheValue := range c {
		if cacheValue.ExpiredAt(now) {
			delete(c, key)
		}
	}
}

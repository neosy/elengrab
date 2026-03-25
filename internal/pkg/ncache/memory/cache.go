package nmemory

import (
	"time"

	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

// Cache is a generic in-memory cache that maps keys of type K to Items of type T.
// It uses a copier function to create safe copies of values before storing them.
type Cache[K comparable, T any] struct {
	// cache holds the internal mapping from keys to items.
	cache map[K]*Item[T]

	// copier is a function used to copy values before storing them in the cache.
	copier CacheCopier[T]
}

// NewCache creates a new Cache instance with the provided copier function.
//
// The returned Cache uses the supplied copier to create copies of values when
// necessary (for example, during Get operations when returning a value to the
// caller without exposing internal references).
//
// Parameters:
//   - K: the type of keys (must be comparable)
//   - T: the type of values stored in the cache
//   - copier: a function that takes a pointer to T and returns a pointer
//     to a deep copy of that value
//
// The internal map is initialized with make() and will grow as needed.
// The Cache does not perform any automatic eviction or size limiting.
//
// Example usage:
//
//	copier := func(p *User) *User {
//	    cp := *p
//	    return &cp
//	}
//	cache := NewCache[string, User](copier)
func NewCache[K comparable, T any](copier CacheCopier[T]) Cache[K, T] {
	return Cache[K, T]{
		cache:  make(map[K]*Item[T]),
		copier: copier,
	}
}

// NewCacheWithDefaultCopier creates a new Cache that uses the default copier
// implementation (DefaultCopier), which relies on the Copy() method defined
// on the pointer type *T.
//
// This is a convenience constructor for the most common case where the type T
// provides a pointer receiver method Copy() *T. It eliminates the need to
// explicitly pass a copier function.
//
// Type parameters:
//   - K:   the type of the cache keys (must be comparable)
//   - T:   the type of the cached values
//   - PT:  the pointer type (*T or a type with underlying type *T) that
//     satisfies the copyable constraint
//
// Example usage:
//
//	cache := NewCacheWithDefaultCopier[uint64, User, *User]()
//
// Note: The type *T (or equivalent) must implement the Copy() method.
//
//	The implementation should handle nil pointers appropriately.
func NewCacheWithDeaultCopier[K comparable, T any, PT copyable[T]]() Cache[K, T] {
	return NewCache[K, T](DefaultCopier[T, PT]())
}

// Save stores a value in the cache with an optional copy function and TTL.
// If ttl > 0, the item expires after the specified duration.
func (c *Cache[K, T]) Save(key K, value *T, ttl time.Duration) {
	var valueCopy *T
	if value != nil {
		valueCopy = uptr.Copy(value)
		if c.copier != nil {
			valueCopy = c.copier(value)
		}
	}

	expiresAt := time.Time{}
	if ttl > 0 {
		expiresAt = time.Now().UTC().Add(ttl)
	}

	var cacheStatus = CacheStatusHit
	if valueCopy == nil {
		cacheStatus = CacheStatusNegativeHit
	}

	c.cache[key] = &Item[T]{
		value:     valueCopy,
		status:    cacheStatus,
		expiresAt: expiresAt,
	}
}

// Delete removes the item with the given key from the cache
func (c Cache[K, T]) Delete(key K) {
	if _, exists := c.cache[key]; !exists {
		return
	}
	delete(c.cache, key)
}

// FindWithStatus returns the cached value for the given key, or nil if not found or expired.
// An optional copy function can be used to return a copy of the value.
func (c Cache[K, T]) FindWithStatus(key K) (*T, CacheStatus) {
	cacheData, exists := c.cache[key]
	if !exists {
		return nil, CacheStatusMiss
	}

	if cacheData.Expired() {
		delete(c.cache, key)
		return nil, CacheStatusMiss
	}

	valueCopy := cacheData.value
	if c.copier != nil {
		valueCopy = c.copier(cacheData.value)
	}

	return valueCopy, cacheData.status
}

// Find returns the cached value for the given key, or nil if not found or expired.
// An optional copy function can be used to return a copy of the value.
func (c Cache[K, T]) Find(key K) *T {
	cacheData, exists := c.cache[key]
	if !exists {
		return nil
	}

	if cacheData.Expired() {
		delete(c.cache, key)
		return nil
	}

	valueCopy := cacheData.value
	if c.copier != nil {
		valueCopy = c.copier(cacheData.value)
	}

	return valueCopy
}

// Exists checks if the cache contains a non-expired item for the given key.
func (c Cache[K, T]) Exists(key K) bool {
	_, status := c.FindWithStatus(key)
	if status == CacheStatusHit {
		return true
	}
	return false
}

// CleanExpired removes all expired items from the cache.
func (c Cache[K, T]) CleanExpired() {
	now := time.Now().UTC()
	for key, cacheValue := range c.cache {
		if cacheValue.ExpiredAt(now) {
			delete(c.cache, key)
		}
	}
}

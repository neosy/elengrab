package memsimple

import (
	"time"

	iptr "github.com/neosy/elengrab/internal/pkg/cache/internal/pointer"
)

// Cache is a generic in-memory cache that maps keys of type K to Items of type T.
// It uses a copier function to create safe copies of values before storing them.
type Cache[K comparable, T any] struct {
	// cache holds the internal mapping from keys to items.
	cache map[K]*Item[T]

	// copier is a function used to copy values before storing them in the cache.
	copier CacheCopier[T]

	// NegativeTTL is a special duration value that indicates a negative cache entry.
	// When a lookup results in a cache miss, the cache can store a negative entry
	// with this TTL to prevent repeated lookups for the same missing key.
	NegativeTTL NegativeTTL
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
//	cache := NewCacheWithDeaultCopier[uint64, User, *User]()
//
// Note: The type *T (or equivalent) must implement the Copy() method.
//
//	The implementation should handle nil pointers appropriately.
func NewCacheWithDeaultCopier[K comparable, T any, PT copyable[T]]() Cache[K, T] {
	return NewCache[K](DefaultCopier[T, PT]())
}

// negativeTTL returns the NegativeTTL function to use for this cache.
func (c *Cache[K, T]) negativeTTL() NegativeTTL {
	if c.NegativeTTL != nil {
		return c.NegativeTTL
	}
	return DefaultNegativeTTL
}

// Save stores a value in the cache with an optional TTL.
// If the caller wants maximum performance, they should use SaveWithNow
// and pass the current time themselves to avoid repeated calls to time.Now().
func (c *Cache[K, T]) Save(key K, value *T, ttl time.Duration) {
	c.SaveWithNow(key, value, ttl, time.Now().UTC())
}

// SaveWithNow is the high-performance version of Save.
// The caller provides a precomputed current time to avoid repeated time.Now()
// calls in hot paths such as batch inserts or bulk cache updates.
func (c *Cache[K, T]) SaveWithNow(key K, value *T, ttl time.Duration, now time.Time) {
	var valueCopy *T
	if value != nil {
		if c.copier == nil {
			valueCopy = iptr.Copy(value)
		} else {
			valueCopy = c.copier(value)
		}
	}

	expiresAt := time.Time{}
	if ttl > 0 {
		if valueCopy == nil {
			expiresAt = now.Add(c.negativeTTL()(ttl))
		} else {
			expiresAt = now.Add(ttl)
		}
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
func (c *Cache[K, T]) Delete(key K) {
	delete(c.cache, key)
}

// FindWithStatus returns the cached value for the given key, or nil if not found or expired.
// An optional copy function can be used to return a copy of the value.
func (c *Cache[K, T]) FindWithStatus(key K) (*T, CacheStatus) {
	return c.FindWithStatusNow(key, time.Now().UTC())
}

// FindWithStatusNow is the high-performance version of FindWithStatus.
// The caller provides the current time to avoid repeated time.Now() calls
// in hot paths such as batch lookups or high-load cache access.
func (c *Cache[K, T]) FindWithStatusNow(key K, now time.Time) (*T, CacheStatus) {
	cacheData, exists := c.cache[key]
	if !exists {
		return nil, CacheStatusMiss
	}

	if cacheData.ExpiredWithNow(now) {
		return nil, CacheStatusMiss
	}

	valueCopy := cacheData.value
	if c.copier != nil {
		valueCopy = c.copier(cacheData.value)
	}

	return valueCopy, cacheData.status
}

// Find returns the cached value for the given key, or nil if not found or expired.
// This is a convenience method that calls FindWithStatus and ignores the status.
func (c *Cache[K, T]) Find(key K) *T {
	return c.FindWithNow(key, time.Now().UTC())
}

// Find returns the cached value for the given key, or nil if not found or expired.
// An optional copy function can be used to return a copy of the value.
func (c *Cache[K, T]) FindWithNow(key K, now time.Time) *T {
	cacheData, exists := c.cache[key]
	if !exists {
		return nil
	}

	if cacheData.ExpiredWithNow(now) {
		return nil
	}

	valueCopy := cacheData.value
	if c.copier != nil {
		valueCopy = c.copier(cacheData.value)
	}

	return valueCopy
}

// ExistsWithNow checks if a key exists in the cache and is not expired.
// The caller provides the current time to avoid repeated time.Now() calls
// in hot paths such as batch lookups or high-load cache access.
func (c *Cache[K, T]) ExistsWithNow(key K, now time.Time) bool {
	cacheData, exists := c.cache[key]
	if !exists {
		return false
	}

	if cacheData.ExpiredWithNow(now) {
		return false
	}

	return true
}

// ExistsWithStatus checks if a key exists in the cache and is not expired, returning its status.
func (c *Cache[K, T]) ExistsWithStatus(key K) (bool, CacheStatus) {
	return c.ExistsWithStatusNow(key, time.Now().UTC())
}

// ExistsWithStatusNow checks if a key exists in the cache and is not expired, returning its status.
// The caller provides the current time to avoid repeated time.Now() calls
// in hot paths such as batch lookups or high-load cache access.
func (c *Cache[K, T]) ExistsWithStatusNow(key K, now time.Time) (bool, CacheStatus) {
	cacheData, exists := c.cache[key]
	if !exists {
		return false, CacheStatusMiss
	}

	if cacheData.ExpiredWithNow(now) {
		return false, CacheStatusMiss
	}

	return true, cacheData.status
}

// Exists checks if a key exists in the cache and is not expired.
// This is a convenience method that calls ExistsWithNow and ignores the current time.
func (c *Cache[K, T]) Exists(key K) bool {
	return c.ExistsWithNow(key, time.Now().UTC())
}

// CleanExpired removes all expired items from the cache.
func (c *Cache[K, T]) CleanExpired() {
	c.CleanExpiredWithNow(time.Now().UTC())
}

// CleanExpiredWithNow is the high-performance version of CleanExpired.
// The caller provides the current time to eliminate the cost of time.Now().UTC()
// in the hot path. Use this when doing batch cleanups or under high load.
func (c *Cache[K, T]) CleanExpiredWithNow(now time.Time) {
	for key, cacheValue := range c.cache {
		if cacheValue.ExpiredAt(now) {
			delete(c.cache, key)
		}
	}
}

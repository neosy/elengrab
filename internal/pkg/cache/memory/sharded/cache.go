package memsharded

import (
	"fmt"
	"hash/fnv"
	"time"

	iptr "github.com/neosy/elengrab/internal/pkg/cache/internal/pointer"
)

// Cache is a generic in-memory cache that maps keys of type K to Items of type T.
// It uses a cloner function to create safe copies of values before storing them.
type Cache[K comparable, T any] struct {
	// shards holds the internal sharded mapping from keys to items.
	shards []*shard[K, T]

	// count is the number of shards used for distribution.
	count uint64

	// cloner is a function used to copy values before storing them in the cache.
	cloner CacheCloner[T]

	// NegativeTTL is a special duration value that indicates a negative cache entry.
	// When a lookup results in a cache miss, the cache can store a negative entry
	// with this TTL to prevent repeated lookups for the same missing key.
	NegativeTTL NegativeTTL
}

// NewCache creates a new Cache instance with the provided cloner function.
//
// The returned Cache uses the supplied cloner to create copies of values when
// necessary (for example, during Get operations when returning a value to the
// caller without exposing internal references).
//
// Parameters:
//   - K: the type of keys (must be comparable)
//   - T: the type of values stored in the cache
//   - cloner: a function that takes a pointer to T and returns a pointer
//     to a deep copy of that value
//
// The internal shards are initialized with make() and will grow as needed.
// The Cache does not perform any automatic eviction or size limiting.
//
// Example usage:
//
//	cloner := func(p *User) *User {
//	    cp := *p
//	    return &cp
//	}
//	cache := NewCache[string, User](cloner)
func NewCache[K comparable, T any](shardCount int, cloner CacheCloner[T]) Cache[K, T] {
	if shardCount <= 0 {
		shardCount = defaultShardCount
	}

	shards := make([]*shard[K, T], shardCount)
	for i := 0; i < shardCount; i++ {
		shards[i] = newShard[K, T]()
	}

	return Cache[K, T]{
		shards: shards,
		count:  uint64(shardCount),
		cloner: cloner,
	}
}

// NewCacheWithDeaultCloner creates a new Cache that uses the default cloner
// implementation (DefaultCloner), which relies on the Clone() method defined
// on the pointer type *T.
//
// This is a convenience constructor for the most common case where the type T
// provides a pointer receiver method Clone() *T. It eliminates the need to
// explicitly pass a cloner function.
//
// Type parameters:
//   - K:   the type of the cache keys (must be comparable)
//   - T:   the type of the cached values
//   - PT:  the pointer type (*T or a type with underlying type *T) that
//     satisfies the clonable constraint
//
// Example usage:
//
//	cache := NewCacheWithDeaultCloner[uint64, User, *User]()
//
// Note: The type *T (or equivalent) must implement the Clone() method.
//
//	The implementation should handle nil pointers appropriately.
func NewCacheWithDeaultCloner[K comparable, T any, PT clonable[T]]() Cache[K, T] {
	return NewCache[K](defaultShardCount, DefaultCloner[T, PT]())
}

// negativeTTL returns the NegativeTTL function to use for this cache.
func (c *Cache[K, T]) negativeTTL() NegativeTTL {
	if c.NegativeTTL != nil {
		return c.NegativeTTL
	}
	return DefaultNegativeTTL
}

// getShard returns shard for a given key.
func (c *Cache[K, T]) getShard(key K) *shard[K, T] {
	h := hashKey(key)
	return c.shards[h%c.count]
}

// Save stores a value in the cache with an optional TTL.
// If the caller wants maximum performance, they should use SaveWithNow
// and pass the current time themselves to avoid repeated calls to time.Now().
func (c *Cache[K, T]) Save(key K, value *T, ttl time.Duration) {
	c.SaveWithNow(key, value, ttl, time.Now().UTC())
}

// SaveWithNow is the high-performance version of Save.
// The caller provides the current time to eliminate the cost of time.Now().UTC()
// in the hot path. Use this when doing batch inserts or under high load.
func (c *Cache[K, T]) SaveWithNow(key K, value *T, ttl time.Duration, now time.Time) {
	s := c.getShard(key)

	var valueClone *T
	if value != nil {
		if c.cloner == nil {
			valueClone = iptr.Clone(value)
		} else {
			valueClone = c.cloner(value)
		}
	}

	expiresAt := time.Time{}
	if ttl > 0 {
		if valueClone == nil {
			expiresAt = now.Add(c.negativeTTL()(ttl))
		} else {
			expiresAt = now.Add(ttl)
		}
	}

	var cacheStatus = CacheStatusHit
	if valueClone == nil {
		cacheStatus = CacheStatusNegativeHit
	}

	s.cache[key] = &Item[T]{
		value:     valueClone,
		status:    cacheStatus,
		expiresAt: expiresAt,
	}
}

// Delete removes the item with the given key from the cache
func (c *Cache[K, T]) Delete(key K) {
	s := c.getShard(key)
	delete(s.cache, key)
}

// FindWithStatus returns the cached value for the given key, or nil if not found or expired.
// An optional clone function can be used to return a copy of the value.
func (c *Cache[K, T]) FindWithStatus(key K) (*T, CacheStatus) {
	return c.FindWithStatusNow(key, time.Now().UTC())
}

// FindWithStatusWithNow is the high-performance version of FindWithStatus.
// The caller provides the current time to eliminate the cost of time.Now().UTC()
// in the hot path. Use this when doing batch lookups or under high load.
func (c *Cache[K, T]) FindWithStatusNow(key K, now time.Time) (*T, CacheStatus) {
	s := c.getShard(key)

	item, exists := s.cache[key]

	if !exists {
		return nil, CacheStatusMiss
	}

	if item.ExpiredWithNow(now) {
		return nil, CacheStatusMiss
	}

	valueClone := item.value
	if c.cloner != nil && valueClone != nil {
		valueClone = c.cloner(item.value)
	}

	return valueClone, item.status
}

// Find returns the cached value for the given key, or nil if not found or expired.
// This is a convenience method that calls FindWithStatus and ignores the status.
func (c *Cache[K, T]) Find(key K) *T {
	return c.FindWithNow(key, time.Now().UTC())
}

// FindWithNow is the high-performance version of Find.
// The caller provides the current time to eliminate the cost of time.Now().UTC()
// in the hot path. Use this when doing batch lookups or under high load.
func (c *Cache[K, T]) FindWithNow(key K, now time.Time) *T {
	v, _ := c.FindWithStatusNow(key, now)
	return v
}

// Exists checks if the cache contains a non-expired item for the given key.
func (c *Cache[K, T]) Exists(key K) bool {
	_, status := c.FindWithStatus(key)
	return status == CacheStatusHit
}

// CleanExpired removes all expired items from the cache.
func (c *Cache[K, T]) CleanExpired() {
	now := time.Now().UTC()

	for _, s := range c.shards {
		for key, item := range s.cache {
			if item.ExpiredAt(now) {
				delete(s.cache, key)
			}
		}
	}
}

// hashKey returns a hash value for shard selection.
func hashKey[K comparable](key K) uint64 {
	switch v := any(key).(type) {
	case uint64:
		return v ^ (v >> 32) ^ (v >> 16) // simple fast mixing
	case int64:
		u := uint64(v)
		return u ^ (u >> 32) ^ (u >> 16)
	case uint32:
		u := uint64(v)
		return u ^ (u << 13) ^ (u >> 7)
	case int:
		u := uint64(v)
		return u ^ (u >> 32) ^ (u >> 16)
	case string:
		h := fnv.New64a()
		h.Write([]byte(v))
		return h.Sum64()
	default:
		// fallback for rare types
		h := fnv.New64a()
		h.Write(fmt.Append(nil, v))
		return h.Sum64()
	}
}

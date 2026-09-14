package memsharded

import (
	"testing"
	"time"
)

type benchUser struct {
	ID   uint64
	Name string
	Age  int
}

func (u *benchUser) Clone() *benchUser {
	if u == nil {
		return nil
	}
	cp := *u
	return &cp
}

// BenchmarkCache_Save measures Save performance using the compatibility method.
func BenchmarkCache_Save(b *testing.B) {
	cache := NewCacheWithDeaultCloner[uint64, benchUser, *benchUser]()

	user := &benchUser{ID: 1, Name: "Benchmark User", Age: 30}

	b.ResetTimer()
	for i := range b.N {
		cache.Save(uint64(i), user, 10*time.Minute)
	}
}

// BenchmarkCache_SaveWithNow measures the high-performance save path
// when current time is provided by the caller (recommended for hot paths).
func BenchmarkCache_SaveWithNow(b *testing.B) {
	cache := NewCacheWithDeaultCloner[uint64, benchUser, *benchUser]()

	user := &benchUser{ID: 1, Name: "Benchmark User", Age: 30}
	now := time.Now().UTC().UTC()

	b.ResetTimer()
	for i := range b.N {
		cache.SaveWithNow(uint64(i), user, 10*time.Minute, now)
	}
}

// BenchmarkCache_SaveNoClone measures SaveWithNow without any value cloning
// (cloner = nil). This shows the "base" speed of the cache.
func BenchmarkCache_SaveNoClone(b *testing.B) {
	cache := NewCache[uint64, benchUser](64, nil) // explicit nil cloner

	user := &benchUser{ID: 1, Name: "Benchmark User", Age: 30}
	now := time.Now().UTC()

	b.ResetTimer()
	for i := range b.N {
		cache.SaveWithNow(uint64(i), user, 10*time.Minute, now)
	}
}

// BenchmarkCache_FindHit measures read performance on cache hits.
func BenchmarkCache_FindHit(b *testing.B) {
	cache := NewCacheWithDeaultCloner[uint64, benchUser, *benchUser]()

	now := time.Now().UTC()

	// Pre-populate cache
	for i := range 100_000 {
		cache.SaveWithNow(uint64(i), &benchUser{ID: uint64(i)}, 30*time.Minute, now)
	}

	b.ResetTimer()
	for i := range b.N {
		_ = cache.Find(uint64(i % 100_000))
	}
}

func BenchmarkCache_FindHitWithNow(b *testing.B) {
	cache := NewCacheWithDeaultCloner[uint64, benchUser, *benchUser]()

	now := time.Now().UTC()

	// Pre-populate cache
	for i := range 100_000 {
		cache.SaveWithNow(uint64(i), &benchUser{ID: uint64(i)}, 30*time.Minute, now)
	}

	now = time.Now().UTC()

	b.ResetTimer()
	for i := range b.N {
		_ = cache.FindWithNow(uint64(i%100_000), now)
	}
}

// BenchmarkCache_FindWithStatusHit measures FindWithStatus on hits.
func BenchmarkCache_FindWithStatusHit(b *testing.B) {
	cache := NewCacheWithDeaultCloner[uint64, benchUser, *benchUser]()

	for i := 0; i < 100_000; i++ {
		cache.SaveWithNow(uint64(i), &benchUser{ID: uint64(i)}, 30*time.Minute, time.Now().UTC())
	}

	b.ResetTimer()
	for i := range b.N {
		_, _ = cache.FindWithStatus(uint64(i % 100_000))
	}
}

// BenchmarkCache_FindMiss measures performance on cache misses.
func BenchmarkCache_FindMiss(b *testing.B) {
	cache := NewCacheWithDeaultCloner[uint64, benchUser, *benchUser]()

	b.ResetTimer()
	for i := range b.N {
		_, _ = cache.FindWithStatus(uint64(i))
	}
}

// BenchmarkCache_Exists measures Exists check performance.
func BenchmarkCache_Exists(b *testing.B) {
	cache := NewCacheWithDeaultCloner[uint64, benchUser, *benchUser]()

	for i := 0; i < 50_000; i++ {
		cache.SaveWithNow(uint64(i), &benchUser{ID: uint64(i)}, 30*time.Minute, time.Now().UTC())
	}

	b.ResetTimer()
	for i := range b.N {
		_ = cache.Exists(uint64(i % 50_000))
	}
}

// BenchmarkCache_CleanExpired measures cleanup performance.
func BenchmarkCache_CleanExpired(b *testing.B) {
	cache := NewCacheWithDeaultCloner[uint64, benchUser, *benchUser]()

	now := time.Now().UTC()

	// Add mix of valid and expired items
	for i := 0; i < 50_000; i++ {
		cache.SaveWithNow(uint64(i), &benchUser{ID: uint64(i)}, 1*time.Hour, now)
	}
	for i := 50_000; i < 100_000; i++ {
		// Simulate expired items
		item := &Item[benchUser]{
			value:     &benchUser{ID: uint64(i)},
			expiresAt: now.Add(-10 * time.Minute),
			status:    CacheStatusHit,
		}
		s := cache.getShard(uint64(i))
		s.mu.Lock()
		s.cache[uint64(i)] = item
		s.mu.Unlock()
	}

	b.ResetTimer()
	for b.Loop() {
		cache.CleanExpired()
	}
}

// BenchmarkCache_MixedWorkload simulates a realistic 90% read / 10% write workload
// under parallel execution.

type benchRepository struct {
	Repository[benchUser]

	cache Cache[uint64, benchUser]
}

func newBenchRepository() *benchRepository {
	r := &benchRepository{
		cache: NewCacheWithDeaultCloner[uint64, benchUser, *benchUser](),
	}

	r.Repository.Init(30 * time.Minute)

	return r
}

// BenchmarkCache_MixedWorkload simulates a realistic 90% read / 10% write workload
// under parallel execution.
func BenchmarkCache_MixedWorkload(b *testing.B) {
	repo := newBenchRepository()

	for i := range 100_000 {
		key := uint64(i)
		_ = repo.Save(func() error {
			repo.cache.Save(
				key,
				&benchUser{ID: key},
				repo.TTL(),
			)

			return nil
		})
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		counter := uint64(0)
		for pb.Next() {
			key := counter % 100_000
			if counter%10 == 0 {
				// 10% writes
				_ = repo.Save(func() error {
					repo.cache.Save(
						key,
						&benchUser{ID: key},
						repo.TTL(),
					)
					return nil
				})
			} else {
				// 90% reads
				_, _ = repo.Find(func() (*benchUser, error) {
					return repo.cache.Find(key), nil
				})
			}
			counter++
		}
	})
}

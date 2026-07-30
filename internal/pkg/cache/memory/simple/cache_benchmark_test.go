package memsimple

import (
	"testing"
	"time"
)

type benchUser struct {
	ID   uint64
	Name string
	Age  int
}

func (u *benchUser) Copy() *benchUser {
	if u == nil {
		return nil
	}
	cp := *u
	return &cp
}

// BenchmarkCache_Save measures the performance of saving items to the cache.
func BenchmarkCache_Save(b *testing.B) {
	cache := NewCacheWithDeaultCopier[uint64, benchUser, *benchUser]()

	user := &benchUser{ID: 1, Name: "Benchmark User", Age: 30}

	b.ResetTimer()
	for i := range b.N {
		cache.Save(uint64(i), user, 10*time.Minute)
	}
}

func BenchmarkCache_SaveWithNow(b *testing.B) {
	cache := NewCacheWithDeaultCopier[uint64, benchUser, *benchUser]()

	user := &benchUser{ID: 1, Name: "Benchmark User", Age: 30}

	now := time.Now().UTC()

	b.ResetTimer()
	for i := range b.N {
		cache.SaveWithNow(uint64(i), user, 10*time.Minute, now)
	}
}

// BenchmarkCache_Find measures the performance of finding existing (hit) items.
func BenchmarkCache_FindHit(b *testing.B) {
	cache := NewCacheWithDeaultCopier[uint64, benchUser, *benchUser]()

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

// BenchmarkCache_FindMiss measures the performance of cache misses.
func BenchmarkCache_FindMiss(b *testing.B) {
	cache := NewCacheWithDeaultCopier[uint64, benchUser, *benchUser]()

	b.ResetTimer()
	for i := range b.N {
		_ = cache.Find(uint64(i))
	}
}

func BenchmarkCache_FindMissWithNow(b *testing.B) {
	cache := NewCacheWithDeaultCopier[uint64, benchUser, *benchUser]()

	now := time.Now().UTC()

	b.ResetTimer()
	for i := range b.N {
		_ = cache.FindWithNow(uint64(i), now)
	}
}

// BenchmarkCache_FindWithStatus measures the performance of FindWithStatus on hits.
func BenchmarkCache_FindWithStatus(b *testing.B) {
	cache := NewCacheWithDeaultCopier[uint64, benchUser, *benchUser]()

	for i := range 10000 {
		cache.Save(uint64(i), &benchUser{ID: uint64(i)}, 30*time.Minute)
	}

	b.ResetTimer()
	for i := range b.N {
		_, _ = cache.FindWithStatus(uint64(i % 10000))
	}
}

func BenchmarkCache_FindWithStatusNow(b *testing.B) {
	cache := NewCacheWithDeaultCopier[uint64, benchUser, *benchUser]()

	now := time.Now().UTC()

	for i := range 10000 {
		cache.SaveWithNow(uint64(i), &benchUser{ID: uint64(i)}, 30*time.Minute, now)
	}

	now = time.Now().UTC()

	b.ResetTimer()
	for i := range b.N {
		_, _ = cache.FindWithStatusNow(uint64(i%10000), now)
	}
}

// BenchmarkCache_Exists measures the performance of existence checks.
func BenchmarkCache_Exists(b *testing.B) {
	cache := NewCacheWithDeaultCopier[uint64, benchUser, *benchUser]()

	for i := range 10000 {
		cache.Save(uint64(i), &benchUser{ID: uint64(i)}, 30*time.Minute)
	}

	b.ResetTimer()
	for i := range b.N {
		_ = cache.Exists(uint64(i % 10000))
	}
}

// BenchmarkCache_CleanExpired measures the performance of cleaning expired items.
func BenchmarkCache_CleanExpired(b *testing.B) {
	cache := NewCacheWithDeaultCopier[uint64, benchUser, *benchUser]()

	baseNow := time.Now().UTC()

	for i := 0; i < 5000; i++ {
		cache.SaveWithNow(uint64(i), &benchUser{ID: uint64(i)}, 1*time.Hour, baseNow)
	}
	for i := 5000; i < 10000; i++ {
		item := &Item[benchUser]{
			value:     &benchUser{ID: uint64(i)},
			expiresAt: baseNow.Add(-10 * time.Minute),
			status:    CacheStatusHit,
		}
		cache.cache[uint64(i)] = item
	}

	cleanNow := time.Now().UTC()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.CleanExpiredWithNow(cleanNow)
	}
}

// BenchmarkCache_MixedWorkload simulates a realistic mixed read/write workload.
type benchRepository struct {
	Repository[benchUser]

	cache Cache[uint64, benchUser]
}

func newBenchRepository() *benchRepository {
	r := &benchRepository{
		cache: NewCacheWithDeaultCopier[uint64, benchUser, *benchUser](),
	}

	r.Repository.Init(30 * time.Minute)

	return r
}

func BenchmarkRepository_MixedWorkload(b *testing.B) {
	repo := newBenchRepository()

	for i := range 10000 {
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

			key := counter % 10000

			if counter%10 == 0 {

				_ = repo.Save(func() error {
					repo.cache.Save(
						key,
						&benchUser{ID: key},
						repo.TTL(),
					)

					return nil
				})

			} else {

				_, _ = repo.Find(func() (*benchUser, error) {
					return repo.cache.Find(key), nil
				})
			}

			counter++
		}
	})
}

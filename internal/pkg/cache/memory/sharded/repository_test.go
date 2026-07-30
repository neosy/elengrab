package memsharded

import (
	"strconv"
	"sync"
	"testing"
	"time"
)

func newTestRepo[T any](ttl time.Duration) (*Repository[T], Cache[int, T]) {
	repo := &Repository[T]{}
	repo.Init(ttl)

	return repo, NewCache[int, T](16, nil)
}

// TestSaveAndFind ensures basic Save/Find correctness.
func TestSaveAndFind(t *testing.T) {
	repo, cache := newTestRepo[string](time.Minute)

	val := "hello"

	if err := repo.Save(func() error {
		cache.Save(1, &val, repo.TTL())
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	got, err := repo.Find(func() (*string, error) {
		return cache.Find(1), nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if got == nil || *got != "hello" {
		t.Fatalf("expected hello, got %+v", got)
	}
}

// TestExists ensures Exists works for present and absent keys.
func TestExists(t *testing.T) {
	repo, cache := newTestRepo[string](time.Minute)

	val := "world"

	if err := repo.Save(func() error {
		cache.Save(2, &val, repo.TTL())
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	exists, err := repo.Exists(func() (bool, error) {
		return cache.Exists(2), nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if !exists {
		t.Fatal("expected key to exist")
	}

	exists, err = repo.Exists(func() (bool, error) {
		return cache.Exists(999), nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if exists {
		t.Fatal("expected key to be absent")
	}
}

// TestDelete ensures keys are properly removed.
func TestDelete(t *testing.T) {
	repo, cache := newTestRepo[string](time.Minute)

	val := "delete-me"

	_ = repo.Save(func() error {
		cache.Save(3, &val, repo.TTL())
		return nil
	})

	if err := repo.Delete(func() error {
		cache.Delete(3)
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	exists, err := repo.Exists(func() (bool, error) {
		return cache.Exists(3), nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if exists {
		t.Fatal("expected key to be deleted")
	}
}

// TestTTLExpiration ensures TTL correctly expires entries.
func TestTTLExpiration(t *testing.T) {
	repo, cache := newTestRepo[string](10 * time.Millisecond)

	val := "expire"

	_ = repo.Save(func() error {
		cache.Save(4, &val, repo.TTL())
		return nil
	})

	time.Sleep(30 * time.Millisecond)

	got, err := repo.Find(func() (*string, error) {
		return cache.Find(4), nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if got != nil {
		t.Fatalf("expected expired value to be nil")
	}
}

// TestNegativeCache ensures nil values are stored as negative hits.
func TestNegativeCache(t *testing.T) {
	repo, cache := newTestRepo[string](time.Minute)

	_ = repo.Save(func() error {
		cache.Save(5, nil, repo.TTL())
		return nil
	})

	val, status, err := repo.FindWithStatus(func() (*string, CacheStatus, error) {
		v, s := cache.FindWithStatus(5)
		return v, s, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if val != nil {
		t.Fatal("expected nil value")
	}

	if status != CacheStatusNegativeHit {
		t.Fatalf("expected NegativeHit, got %v", status)
	}
}

// TestShardDistribution ensures keys are distributed across shards.
func TestShardDistribution(t *testing.T) {
	repo, cache := newTestRepo[string](time.Minute)

	for i := 0; i < 1000; i++ {
		val := "v"

		_ = repo.Save(func() error {
			cache.Save(i, &val, repo.TTL())
			return nil
		})
	}

	for i := 0; i < 1000; i++ {
		v, err := repo.Find(func() (*string, error) {
			return cache.Find(i), nil
		})
		if err != nil {
			t.Fatal(err)
		}

		if v == nil {
			t.Fatalf("expected value for key %d", i)
		}
	}
}

// TestConcurrentWrites ensures no race conditions under concurrent writes.
func TestConcurrentWrites(t *testing.T) {
	repo, cache := newTestRepo[string](time.Minute)

	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			val := "v" + strconv.Itoa(i)

			_ = repo.Save(func() error {
				cache.Save(i, &val, repo.TTL())
				return nil
			})
		}(i)
	}

	wg.Wait()

	for i := 0; i < 1000; i++ {
		v, err := repo.Find(func() (*string, error) {
			return cache.Find(i), nil
		})
		if err != nil {
			t.Fatal(err)
		}

		if v == nil {
			t.Fatalf("expected value for key %d", i)
		}
	}
}

// TestCopier ensures copier function is applied.
func TestCopier(t *testing.T) {
	type User struct {
		Name string
	}

	copier := func(u *User) *User {
		if u == nil {
			return nil
		}

		cp := *u
		cp.Name = "copied"

		return &cp
	}

	repo := &Repository[User]{}
	repo.Init(time.Minute)

	cache := NewCache[int, User](16, copier)

	u := User{Name: "original"}

	_ = repo.Save(func() error {
		cache.Save(1, &u, repo.TTL())
		return nil
	})

	got, err := repo.Find(func() (*User, error) {
		return cache.Find(1), nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if got == nil || got.Name != "copied" {
		t.Fatalf("expected copied value, got %+v", got)
	}
}

// BenchmarkSave measures write performance.
func BenchmarkSave(b *testing.B) {
	repo, cache := newTestRepo[string](time.Minute)

	val := "bench"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = repo.Save(func() error {
			cache.Save(i, &val, repo.TTL())
			return nil
		})
	}
}

// BenchmarkFind measures read performance.
func BenchmarkFind(b *testing.B) {
	repo, cache := newTestRepo[string](time.Minute)

	val := "bench"

	for i := 0; i < 10000; i++ {
		_ = repo.Save(func() error {
			cache.Save(i, &val, repo.TTL())
			return nil
		})
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = repo.Find(func() (*string, error) {
			return cache.Find(i % 10000), nil
		})
	}
}

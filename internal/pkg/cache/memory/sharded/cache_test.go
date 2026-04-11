package memsharded

import (
	"strconv"
	"sync"
	"testing"
	"time"
)

// TestSaveAndFind ensures basic Save/Find correctness.
func TestSaveAndFind(t *testing.T) {
	cache := NewCache[int, string](16, nil)

	val := "hello"
	cache.Save(1, &val, 0)

	got := cache.Find(1)
	if got == nil || *got != "hello" {
		t.Fatalf("expected hello, got %+v", got)
	}
}

// TestExists ensures Exists works for present and absent keys.
func TestExists(t *testing.T) {
	cache := NewCache[int, string](16, nil)

	val := "world"
	cache.Save(2, &val, 0)

	if !cache.Exists(2) {
		t.Fatalf("expected key to exist")
	}

	if cache.Exists(999) {
		t.Fatalf("expected key to be absent")
	}
}

// TestDelete ensures keys are properly removed.
func TestDelete(t *testing.T) {
	cache := NewCache[int, string](16, nil)

	val := "delete-me"
	cache.Save(3, &val, 0)

	cache.Delete(3)

	if cache.Exists(3) {
		t.Fatalf("expected key to be deleted")
	}
}

// TestTTLExpiration ensures TTL correctly expires entries.
func TestTTLExpiration(t *testing.T) {
	cache := NewCache[int, string](16, nil)

	val := "expire"
	cache.Save(4, &val, 10*time.Millisecond)

	time.Sleep(30 * time.Millisecond)

	got := cache.Find(4)
	if got != nil {
		t.Fatalf("expected expired value to be nil")
	}
}

// TestNegativeCache ensures nil values are stored as negative hits.
func TestNegativeCache(t *testing.T) {
	cache := NewCache[int, string](16, nil)

	cache.Save(5, nil, 0)

	val, status := cache.FindWithStatus(5)
	if val != nil {
		t.Fatalf("expected nil value")
	}

	if status != CacheStatusNegativeHit {
		t.Fatalf("expected NegativeHit, got %v", status)
	}
}

// TestShardDistribution ensures keys are distributed across shards.
func TestShardDistribution(t *testing.T) {
	cache := NewCache[int, string](32, nil)

	// track shard usage indirectly by checking lock contention-free writes
	for i := 0; i < 1000; i++ {
		val := "v"
		cache.Save(i, &val, 0)
	}

	// verify all values are still accessible (correct routing)
	for i := 0; i < 1000; i++ {
		v := cache.Find(i)
		if v == nil {
			t.Fatalf("expected value for key %d", i)
		}
	}
}

// TestConcurrentWrites ensures no race conditions under concurrent writes.
func TestConcurrentWrites(t *testing.T) {
	cache := NewCache[int, string](32, nil)

	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			val := "v" + strconv.Itoa(i)
			cache.Save(i, &val, 0)
		}(i)
	}

	wg.Wait()

	for i := 0; i < 1000; i++ {
		v := cache.Find(i)
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

	cache := NewCache[int, User](16, copier)

	u := User{Name: "original"}
	cache.Save(1, &u, 0)

	got := cache.Find(1)
	if got == nil || got.Name != "copied" {
		t.Fatalf("expected copied value, got %+v", got)
	}
}

// BenchmarkSave measures write performance.
func BenchmarkSave(b *testing.B) {
	cache := NewCache[int, string](64, nil)

	val := "bench"
	for i := 0; i < b.N; i++ {
		cache.Save(i, &val, 0)
	}
}

// BenchmarkFind measures read performance.
func BenchmarkFind(b *testing.B) {
	cache := NewCache[int, string](64, nil)

	val := "bench"
	for i := 0; i < 10000; i++ {
		cache.Save(i, &val, 0)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cache.Find(i % 10000)
	}
}

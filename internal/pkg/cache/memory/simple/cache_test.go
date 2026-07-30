package memsimple

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testUser struct {
	ID   uint64
	Name string
}

func (u *testUser) Copy() *testUser {
	if u == nil {
		return nil
	}
	cp := *u
	return &cp
}

// ====================== CacheStatus ======================

func TestCacheStatus_String(t *testing.T) {
	tests := []struct {
		status   CacheStatus
		expected string
	}{
		{CacheStatusMiss, "miss"},
		{CacheStatusHit, "hit"},
		{CacheStatusNegativeHit, "negative"},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, tt.status.String())
	}
}

func TestCacheStatus_Exists(t *testing.T) {
	assert.True(t, CacheStatusMiss.Exists())
	assert.True(t, CacheStatusHit.Exists())
	assert.True(t, CacheStatusNegativeHit.Exists())

	var invalid CacheStatus = 255
	assert.False(t, invalid.Exists())
}

func TestParseCacheStatus(t *testing.T) {
	tests := []struct {
		input    string
		expected CacheStatus
		hasError bool
	}{
		{"miss", CacheStatusMiss, false},
		{"hit", CacheStatusHit, false},
		{"negative", CacheStatusNegativeHit, false},
		{"MISS", CacheStatusMiss, false},
		{"Hit", CacheStatusHit, false},
		{"NEGATIVE", CacheStatusNegativeHit, false},
		{"invalid", CacheStatusMiss, true},
		{"", CacheStatusMiss, true},
	}

	for _, tt := range tests {
		status, err := ParseCacheStatus(tt.input)
		if tt.hasError {
			assert.Error(t, err)
			assert.Equal(t, CacheStatusMiss, status)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, status)
		}
	}
}

func TestMustParseCacheStatus(t *testing.T) {
	assert.Equal(t, CacheStatusMiss, MustParseCacheStatus("invalid"))
	assert.Equal(t, CacheStatusHit, MustParseCacheStatus("hit"))
}

// ====================== Item ======================

func TestItem_Valid_Expired(t *testing.T) {
	now := time.Now().UTC()

	item := &Item[testUser]{
		value:     &testUser{ID: 1, Name: "Test"},
		expiresAt: time.Time{}, // never expires
	}

	assert.True(t, item.Valid())
	assert.False(t, item.Expired())

	// Expires in future
	item.expiresAt = now.Add(10 * time.Minute)
	assert.True(t, item.Valid())
	assert.False(t, item.Expired())

	// Already expired
	item.expiresAt = now.Add(-10 * time.Minute)
	assert.False(t, item.Valid())
	assert.True(t, item.Expired())
}

func TestItem_ValidAt_ExpiredAt(t *testing.T) {
	now := time.Now().UTC()
	item := &Item[testUser]{
		expiresAt: now.Add(5 * time.Minute),
	}

	assert.True(t, item.ValidAt(now))
	assert.False(t, item.ExpiredAt(now))

	assert.True(t, item.ValidAt(now.Add(5*time.Minute))) // exactly at expiration = still valid
	assert.False(t, item.ExpiredAt(now.Add(5*time.Minute)))

	assert.False(t, item.ValidAt(now.Add(6*time.Minute)))
	assert.True(t, item.ExpiredAt(now.Add(6*time.Minute)))
}

// ====================== Cache ======================

func TestNewCache(t *testing.T) {
	copier := func(u *testUser) *testUser {
		if u == nil {
			return nil
		}
		cp := *u
		return &cp
	}

	cache := NewCache[string, testUser](copier)
	assert.NotNil(t, cache.cache)
	assert.NotNil(t, cache.copier)
}

func TestNewCacheWithDefaultCopier(t *testing.T) {
	cache := NewCacheWithDeaultCopier[uint64, testUser, *testUser]()
	assert.NotNil(t, cache.cache)
	assert.NotNil(t, cache.copier)
}

func TestCache_Save_And_Find(t *testing.T) {
	cache := NewCacheWithDeaultCopier[uint64, testUser, *testUser]()

	user := &testUser{ID: 42, Name: "Alice"}

	// Save with TTL
	cache.Save(1, user, 10*time.Minute)

	// Find should return copy
	found := cache.Find(1)
	require.NotNil(t, found)
	assert.Equal(t, user.ID, found.ID)
	assert.Equal(t, user.Name, found.Name)

	// Should be different pointer (copy)
	assert.NotSame(t, user, found)

	// Negative cache (nil value)
	cache.Save(2, nil, 0)
	found = cache.Find(2)
	assert.Nil(t, found)

	_, status := cache.FindWithStatus(2)
	assert.Equal(t, CacheStatusNegativeHit, status)
}

func TestCache_FindWithStatus(t *testing.T) {
	cache := NewCacheWithDeaultCopier[uint64, testUser, *testUser]()

	cache.Save(1, &testUser{ID: 1}, 0)
	val, status := cache.FindWithStatus(1)
	assert.NotNil(t, val)
	assert.Equal(t, CacheStatusHit, status)

	// Miss
	val, status = cache.FindWithStatus(999)
	assert.Nil(t, val)
	assert.Equal(t, CacheStatusMiss, status)
}

func TestCache_Exists(t *testing.T) {
	cache := NewCacheWithDeaultCopier[uint64, testUser, *testUser]()

	assert.False(t, cache.Exists(1))

	cache.Save(1, &testUser{ID: 1}, 0)
	assert.True(t, cache.Exists(1))

	// Negative hit should return false for Exists
	cache.Save(2, nil, 0)
	assert.False(t, cache.Exists(2))
}

func TestCache_Expiration(t *testing.T) {
	cache := NewCacheWithDeaultCopier[uint64, testUser, *testUser]()

	cache.Save(1, &testUser{ID: 1}, 50*time.Millisecond)

	assert.True(t, cache.Exists(1))

	time.Sleep(100 * time.Millisecond)

	assert.False(t, cache.Exists(1))
	assert.Nil(t, cache.Find(1))
}

func TestCache_CleanExpired(t *testing.T) {
	cache := NewCacheWithDeaultCopier[uint64, testUser, *testUser]()

	cache.Save(1, &testUser{ID: 1}, 0)                   // never expires
	cache.Save(2, &testUser{ID: 2}, 10*time.Millisecond) // will expire
	cache.Save(3, nil, 0)                                // negative

	time.Sleep(50 * time.Millisecond)

	cache.CleanExpired()

	assert.True(t, cache.Exists(1))
	assert.False(t, cache.Exists(2))
	assert.False(t, cache.Exists(3)) // negative also removed on clean? Wait — check logic
}

func TestCache_Delete(t *testing.T) {
	cache := NewCacheWithDeaultCopier[uint64, testUser, *testUser]()

	cache.Save(1, &testUser{ID: 1}, 0)
	assert.True(t, cache.Exists(1))

	cache.Delete(1)
	assert.False(t, cache.Exists(1))
}

// ====================== Repository ======================

func TestRepository(t *testing.T) {
	repo := &Repository[testUser]{}
	repo.Init(5 * time.Minute)

	assert.Equal(t, 5*time.Minute, repo.TTL())

	called := false
	err := repo.Save(func() error {
		called = true
		return nil
	})
	assert.NoError(t, err)
	assert.True(t, called)

	// Test Find
	called = false
	_, err = repo.Find(func() (*testUser, error) {
		called = true
		return nil, nil
	})
	assert.NoError(t, err)
	assert.True(t, called)
}

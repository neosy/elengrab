package memsimple

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testValue struct {
	ID   int
	Name string
}

func TestRepository_Init(t *testing.T) {
	var repo Repository[testValue]

	ttl := 10 * time.Second

	repo.Init(ttl)

	assert.Equal(t, ttl, repo.TTL())
}

func TestRepository_Save(t *testing.T) {
	var repo Repository[testValue]
	repo.Init(time.Second)

	called := false

	err := repo.Save(func() error {
		called = true
		return nil
	})

	require.NoError(t, err)
	assert.True(t, called)
}

func TestRepository_SaveError(t *testing.T) {
	var repo Repository[testValue]
	repo.Init(time.Second)

	expected := errors.New("save error")

	err := repo.Save(func() error {
		return expected
	})

	assert.ErrorIs(t, err, expected)
}

func TestRepository_Delete(t *testing.T) {
	var repo Repository[testValue]
	repo.Init(time.Second)

	called := false

	err := repo.Delete(func() error {
		called = true
		return nil
	})

	require.NoError(t, err)
	assert.True(t, called)
}

func TestRepository_Find(t *testing.T) {
	var repo Repository[testValue]
	repo.Init(time.Second)

	value := testValue{
		ID: 1,
	}

	result, err := repo.Find(func() (*testValue, error) {
		return &value, nil
	})

	require.NoError(t, err)

	require.NotNil(t, result)
	assert.Equal(t, value, *result)
}

func TestRepository_FindError(t *testing.T) {
	var repo Repository[testValue]
	repo.Init(time.Second)

	expected := errors.New("find error")

	_, err := repo.Find(func() (*testValue, error) {
		return nil, expected
	})

	assert.ErrorIs(t, err, expected)
}

func TestRepository_FindWithStatus(t *testing.T) {
	var repo Repository[testValue]
	repo.Init(time.Second)

	value := testValue{
		ID: 10,
	}

	result, status, err := repo.FindWithStatus(func() (*testValue, CacheStatus, error) {
		return &value, CacheStatusHit, nil
	})

	require.NoError(t, err)

	assert.Equal(t, CacheStatusHit, status)
	assert.Equal(t, value, *result)
}

func TestRepository_Exists(t *testing.T) {
	var repo Repository[testValue]
	repo.Init(time.Second)

	exists, err := repo.Exists(func() (bool, error) {
		return true, nil
	})

	require.NoError(t, err)
	assert.True(t, exists)
}

func TestRepository_CleanExpired(t *testing.T) {
	var repo Repository[testValue]
	repo.Init(time.Second)

	called := false

	err := repo.CleanExpired(func() error {
		called = true
		return nil
	})

	require.NoError(t, err)

	assert.True(t, called)
}

func TestRepository_CopyAdapter(t *testing.T) {
	var repo Repository[testValue]

	expected := &testValue{
		ID: 5,
	}

	copier := repo.CopyAdapter(func() *testValue {
		return expected
	})

	result := copier(nil)

	assert.Equal(t, expected, result)
}

func TestRepository_ConcurrentRead(t *testing.T) {
	var repo Repository[testValue]
	repo.Init(time.Second)

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			_, err := repo.Find(func() (*testValue, error) {
				time.Sleep(time.Millisecond)
				return &testValue{}, nil
			})

			require.NoError(t, err)
		}()
	}

	wg.Wait()
}

func TestRepository_ConcurrentWrite(t *testing.T) {
	var repo Repository[testValue]
	repo.Init(time.Second)

	var (
		counter int
		wg      sync.WaitGroup
	)

	for range 1000 {
		wg.Add(1)

		go func() {
			defer wg.Done()

			err := repo.Save(func() error {
				counter++
				return nil
			})

			require.NoError(t, err)
		}()
	}

	wg.Wait()

	assert.Equal(t, 1000, counter)
}

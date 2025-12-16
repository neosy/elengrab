package cache

import (
	"time"
)

type CacheMap[K comparable, T any] map[K]*Entry[T]

func (c CacheMap[K, T]) Save(key K, value *T, fnCopy func(value *T) *T, ttl time.Duration) {
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

	c[key] = &Entry[T]{
		value:     valueCopy,
		expiresAt: expiresAt,
	}
}

func (c CacheMap[K, T]) Delete(key K) {
	if _, exists := c[key]; !exists {
		return
	}
	delete(c, key)
}

func (c CacheMap[K, T]) Find(key K, fnCopy func(value *T) *T) *T {
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

func (c CacheMap[K, T]) Exists(key K) bool {
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

func (c CacheMap[K, T]) CleanExpired() {
	for key, cacheValue := range c {
		if cacheValue.Expired() {
			delete(c, key)
		}
	}
}

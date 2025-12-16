package cache

import "time"

type Entry[T any] struct {
	value     *T
	expiresAt time.Time
}

func (e *Entry[T]) Expired() bool {
	if e.expiresAt.IsZero() {
		return false
	}
	return time.Now().After(e.expiresAt)
}

package nmemory

import "time"

// Item represents a cache entry storing a value and its expiration time.
type Item[T any] struct {
	value     *T
	expiresAt time.Time
}

// Expired returns true if the item has expired.
func (e *Item[T]) Expired() bool {
	if e.expiresAt.IsZero() {
		return false
	}
	return time.Now().After(e.expiresAt)
}

// ExpiredAt returns true if the item has expired at the given time.
func (e *Item[T]) ExpiredAt(now time.Time) bool {
	if e.expiresAt.IsZero() {
		return false
	}
	return now.After(e.expiresAt)
}

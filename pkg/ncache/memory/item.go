package nmemory

import "time"

// Item represents a cache entry storing a value and its expiration time.
type Item[T any] struct {
	value     *T
	expiresAt time.Time
}

// Valid reports whether the item is still valid (not expired).
func (e *Item[T]) Valid() bool {
	return e.expiresAt.IsZero() || !time.Now().After(e.expiresAt)
}

// Expired reports whether the item has expired.
// Returns true only if an expiration time is set and the current time is strictly after it.
func (e *Item[T]) Expired() bool {
	return !e.Valid()
}

// ValidAt reports whether the item is still valid at the given time.
// Returns true if:
//   - the expiration time is zero (never expires), or
//   - now ≤ e.expiresAt (the item is valid at the exact moment of expiration).
func (e *Item[T]) ValidAt(now time.Time) bool {
	return e.expiresAt.IsZero() || !now.After(e.expiresAt)
}

// ExpiredAt reports whether the item has expired at the given time.
// Returns true only if an expiration time is set and now > e.expiresAt.
func (e *Item[T]) ExpiredAt(now time.Time) bool {
	return !e.ValidAt(now)
}

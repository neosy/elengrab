package uptr

// Deref returns the value from the given pointer (dereference),
// or the zero value of T if the pointer is nil.
func Deref[T any](ptr *T) T {
	if ptr != nil {
		return *ptr
	}

	var zero T
	return zero
}

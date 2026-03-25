package uptr

// Copy returns a new pointer to a copy of v, or nil if v is nil.
func Copy[T any](v *T) *T {
	if v == nil {
		return nil
	}

	copy := *v

	return &copy
}

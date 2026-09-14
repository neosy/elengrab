package uptr

// Clone returns a new pointer to a copy of v, or nil if v is nil.
func Clone[T any](v *T) *T {
	if v == nil {
		return nil
	}
	return new(*v)
}

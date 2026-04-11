package memsharded

// CacheCopier is a function type that takes a pointer to a value of type T
// and returns a pointer to a deep copy of that value.
type CacheCopier[T any] func(*T) *T

// copyable is a constraint that requires the implementing type to be
// (or have an underlying type of) a pointer to T and to provide a Copy()
// method that returns a pointer to a copy of the value.
type copyable[T any] interface {
	~*T
	Copy() *T
}

// DefaultCopier returns a CacheCopier function that creates copies using
// the Copy() method defined on pointer types.
//
// Type parameters:
//   - T:   the base (non-pointer) type being copied
//   - PT:  the pointer type (*T or a type with underlying type *T)
//     that implements the copyable constraint
//
// Usage example:
//
//	DefaultCopier[Person, *Person]()
//	DefaultCopier[User, *User]()
func DefaultCopier[T any, PT copyable[T]]() CacheCopier[T] {
	return func(src *T) *T {
		return PT(src).Copy()
	}
}
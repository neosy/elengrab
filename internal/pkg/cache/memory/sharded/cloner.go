package memsharded

// CacheCloner is a function type that takes a pointer to a value of type T
// and returns a pointer to a deep copy of that value.
type CacheCloner[T any] func(*T) *T

// clonable is a constraint that requires the implementing type to be
// (or have an underlying type of) a pointer to T and to provide a Clone()
// method that returns a pointer to a clone of the value.
type clonable[T any] interface {
	~*T
	Clone() *T
}

// DefaultCloner returns a CacheCloner function that creates clones using
// the Clone() method defined on pointer types.
//
// Type parameters:
//   - T:   the base (non-pointer) type being cloned
//   - PT:  the pointer type (*T or a type with underlying type *T)
//     that implements the clonable constraint
//
// Usage example:
//
//	DefaultCloner[Person, *Person]()
//	DefaultCloner[User, *User]()
func DefaultCloner[T any, PT clonable[T]]() CacheCloner[T] {
	return func(src *T) *T {
		return PT(src).Clone()
	}
}

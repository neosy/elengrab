package fnx

// Ternary is a generic function that mimics the behavior of a ternary operator.
// It returns 'a' if 'cond' is true, otherwise it returns 'b'.
// Example usage:
//
//	result := Ternary(condition, valueIfTrue, valueIfFalse)
func Ternary[T comparable](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}

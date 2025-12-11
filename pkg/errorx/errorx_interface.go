package errorx

import "github.com/neosy/elengrab/pkg/errorx/exceptionx"

// Interface type for error handling
type ErrorxInterface interface {
	error

	// Err returns a value of type error
	Err() error

	// Append merges the given errors into the current error.
	// The pointer remains the same; only the internal field err are updated.
	Append(errs ...error) ErrorxInterface

	// Combine returns a new error that combines the current error with the given ones.
	// A new object is returned; the current object remains unchanged.
	Combine(errs ...error) ErrorxInterface

	// UnwrapAll error analysis errors in a slice
	// Duplicates are excluded
	UnwrapAll() []error

	// UnwrapTexts analyzes errors and then extracts text one by one
	UnwrapTexts() *ErrorTexts

	// Message returns a value of type ErrorxMessage
	Message() ErrorxMessage

	// ExceptionType returns a type of exception
	ExceptionType() exceptionx.ExceptionType

	// ExceptionCode returns a code of exception
	ExceptionCode() exceptionx.ExceptionCode
}

package errorx

import "github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"

// Interface type for error handling
type Errorx interface {
	error

	// Err returns a value of type error
	Err() error

	// Append merges the given errors into the current error.
	// The pointer remains the same; only the internal field err are updated.
	Append(errs ...error) Errorx

	// Combine returns a new error that combines the current error with the given ones.
	// A new object is returned; the current object remains unchanged.
	Combine(errs ...error) Errorx

	// UnwrapAll error analysis errors in a slice
	// Duplicates are excluded
	UnwrapAll() []error

	// UnwrapTexts analyzes errors and then extracts text one by one
	UnwrapTexts() *ErrorTexts

	// Message returns a value of type string
	Message() *string

	// Exception returns a exception
	Exception() exceptionx.Exception

	// HttpStatusCodeRaw returns the explicitly set HTTP status code.
	HttpStatusCodeRaw() *int

	// HttpStatusCode returns the HTTP status code, falling back to the exception if needed.
	HttpStatusCode() *int

	// Args returns the arguments used to configure the error.
	Args() []any
}

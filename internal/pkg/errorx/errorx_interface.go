package errorx

import "github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"

// Interface type for error handling
type Errorx interface {
	error

	// Err returns a value of type error
	Err() error

	// Wrap adds the given errors to the current error using the library's
	// wrapping mechanism.
	//
	// The receiver is mutated in place: only the internal error field is updated.
	// The returned Errorx value is the same object as the original (the pointer
	// remains unchanged). Nil errors are ignored.
	Wrap(errs ...error) Errorx

	// WrapAndMerge returns a new error that combines the current error with the given
	// ones using the library's wrapping mechanism.
	//
	// A new Errorx object is always returned; the original object remains
	// unchanged. Nil errors are ignored.
	WrapAndMerge(errs ...error) Errorx

	// Join merges the current underlying error with the provided errors using errors.Join.
	// The existing e.err is included as the first element in the resulting error chain.
	// Returns the same Errorx instance to allow method chaining.
	Join(errs ...error) Errorx

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

	// Copy copying an error including nested
	Copy() error
}

type errorxInternal interface {
	Errorx
	initFromErrorx(err errorxInternal)
	initFromArgs(args ...any)
	setErr(err error)
	setException(exception exceptionx.Exception)
	normalizeToInnerType() errorxInternal
	// Args returns the arguments used to configure the error.
	args() []any
}

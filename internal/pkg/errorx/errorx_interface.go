package errorx

import "github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"

// Interface type for error handling
type Errorx interface {
	error

	// Err returns a value of type error
	Err() error

	// Wrap adds the given errors to the current error in place using the library's wrapping mechanism.
	// The current Errorx object is mutated; the receiver remains the same. Nil errors are ignored.
	Wrap(errs ...error) Errorx

	// WrapNew returns a new Errorx that wraps the current error along with the provided errors.
	// The original error remains unchanged; nil errors are ignored.
	WrapNew(errs ...error) Errorx

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
	Message() string

	// Exception returns a exception
	Exception() exceptionx.Exception

	// HttpStatusCodeRaw returns the explicitly set HTTP status code.
	HttpStatusCodeRaw() int

	// HttpStatusCode returns the HTTP status code, falling back to the exception if needed.
	HttpStatusCode() int

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

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

	// PublicMessage returns the public-facing error message, prioritizing the explicitly set message,
	// then the message from the associated exception,
	// and finally falling back to the standard HTTP status text if available.
	PublicMessage() string

	// HttpStatus returns the explicitly set HTTP status.
	HttpStatus() int

	// PublicHttpStatus returns the HTTP status, prioritizing the explicitly set status,
	// then the status from the associated exception, and finally falling back to 500 if none is set.
	PublicHttpStatus() int

	// PublicHttpStatusText returns the standard HTTP status text corresponding to the public HTTP status.
	PublicHttpStatusText() string

	// Exception returns a exception
	Exception() exceptionx.Exception

	// Copy copying an error including nested
	Copy() error

	// OuterMessage returns the first (outermost) message found in the error chain.
	OuterMessage() string

	// RootMessage returns the first (outermost) message found in the error chain.
	RootMessage() string

	// OuterException returns the first (outermost) exception found in the error chain.
	OuterException() exceptionx.Exception

	// RootException returns the first (outermost) exception found in the error chain.
	RootException() exceptionx.Exception

	// OuterHttpStatus returns the first (outermost) HTTP status found in the error chain.
	OuterHttpStatus() int

	// RootHttpStatus returns the first (outermost) HTTP status found in the error chain.
	RootHttpStatus() int
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

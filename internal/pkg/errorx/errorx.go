package errorx

import (
	"errors"

	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

// Type for error handling
type errorx struct {
	err            error
	message        *string
	exception      exceptionx.Exception
	httpStatusCode *int
}

// newErrx creating an errorx object from text
func newErrx(text string) (errx Errorx) {
	errx = &errorx{
		err: errors.New(text),
	}

	return
}

// newFromErr creating errorx from error
func newFromErr(err error) Errorx {
	errx := &errorx{}

	if err == nil {
		return errx
	}

	if e, ok := err.(Errorx); ok {
		errx.initFromErrorx(e)
	} else {
		errx.err = err
	}

	if errx.Exception() == nil {
		errs := UnwrapAll(err)
		if len(errs) > 1 {
			combErrx := CombineErrors(errs...)
			if e, ok := combErrx.(Errorx); ok {
				errx.exception = e.Exception()
			}
		}
	}

	return errx
}

// New creates a new Errorx instance from the provided text and optional arguments.
//
// Supported argument types and their effects:
//   - exceptionx.Exception: sets the underlying exception
//   - HttpStatusProvider: extracts and sets the HTTP status code
//   - ErrorMessageProvider: overrides the default message
//   - error: will be combined with the current error using CombineErrors
//
// Arguments are processed in order; later values may override earlier ones where applicable.
func New(text string, args ...any) (errx Errorx) {
	errx = newErrx(text)

	errx.(*errorx).initFromArgs(args...)

	return
}

// NewHTTP creates a new Errorx instance with the provided message and HTTP status code.
// Additional optional arguments can be passed to further configure the error.
// The httpStatusCode is automatically added to the arguments as HttpStatusCodeArg.
// Supported argument types and their effects:
//   - exceptionx.Exception: sets the underlying exception
//   - ErrorMessageProvider: overrides the default message
func NewHTTP(text string, httpStatusCode int, args ...any) (errx Errorx) {
	errx = newErrx(text)

	args = append(args, HttpStatusArg(httpStatusCode))
	args = append(args, ErrorMessageArg(text))
	errx.(*errorx).initFromArgs(args...)

	return
}

// NewFromError creating errorx from error and arguments
// args can have types: Exception, ErrorMessageProvider, error
// Supported argument types and their effects:
//   - exceptionx.Exception: sets the underlying exception
//   - HttpStatusProvider: extracts and sets the HTTP status code
//   - ErrorMessageProvider: overrides the default message
func NewFromError(err error, args ...any) Errorx {
	if err == nil {
		return nil
	}

	errx := newFromErr(err)

	errx.(*errorx).initFromArgs(args...)

	return errx
}

// initFromErrorx initialization of fields from errorx
func (errx *errorx) initFromErrorx(err Errorx) {
	errx.err = err.Err()

	errx.initFromArgs(err.Args()...)
}

// initFromArgs initialize errorx from arguments
// args can have types: Exception, ErrorMessageProvider, error
func (errx *errorx) initFromArgs(args ...any) {
	var message *string
	var exception exceptionx.Exception
	var httpStatusCode *int
	var err error

	for _, arg := range args {
		switch v := arg.(type) {
		case exceptionx.DomainException:
			exception = v.NewException()
		case exceptionx.Exception:
			exception = v
		case ErrorMessageProvider:
			if v != nil {
				message = v()
			}
		case HttpStatusProvider:
			if v != nil {
				httpStatusCode = v()
			}
		case error:
			err = CombineErrors(err, v)
		}
	}

	if message != nil {
		errx.message = message
	}

	if exception != nil {
		errx.exception = exception
	}

	if httpStatusCode != nil && *httpStatusCode != 0 {
		errx.httpStatusCode = httpStatusCode
	}

	if err != nil {
		errx.Append(err)
	}
}

func (e *errorx) Args() []any {
	args := make([]any, 0, 3)

	args = append(args, e.Exception())
	if text := e.Message(); text != nil {
		args = append(args, ErrorMessageArg(*text))
	}
	if code := e.httpStatusCode; code != nil {
		args = append(args, HttpStatusArg(*code))
	}

	return args
}

// Error returns the error text
func (e *errorx) Error() (text string) {
	if e.err != nil {
		text = e.err.Error()
	}

	return
}

// Err returns a value of type error
func (e *errorx) Err() error {
	return e.err
}

// Message returns a value of type string
func (e *errorx) Message() *string {
	return e.message
}

// Exception returns a exception
func (e *errorx) Exception() exceptionx.Exception {
	return e.exception
}

// HttpStatusCodeRaw returns the explicitly set HTTP status code.
func (e *errorx) HttpStatusCodeRaw() *int {
	return e.httpStatusCode
}

// HttpStatusCode returns the HTTP status code, falling back to the exception if needed.
func (e *errorx) HttpStatusCode() *int {
	if e.httpStatusCode != nil {
		return e.httpStatusCode
	}

	if e.exception != nil {
		code := e.exception.HttpStatusCode()
		if code != 0 {
			return &code
		}
	}

	return nil
}

// Append merges the given errors into the current error.
// The pointer remains the same; only the internal field err are updated.
func (e *errorx) Append(errs ...error) Errorx {
	e.err = e.Combine(errs...).Err()
	return e
}

// Combine returns a new error that combines the current error with the given ones.
// A new object is returned; the current object remains unchanged.
func (e *errorx) Combine(errs ...error) Errorx {
	return CombineErrors(e, CombineErrors(errs...)).(Errorx)
}

// UnwrapAll error analysis errors in a slice
// Duplicates are excluded
func (e *errorx) UnwrapAll() []error {
	return UnwrapAll(e)
}

// UnwrapTexts analyzes errors and then extracts text one by one
func (e *errorx) UnwrapTexts() *ErrorTexts {
	return NewErrorTexts().AddUnwrapErr(e)
}

// Copy copying an error including nested
func (e *errorx) Copy() Errorx {
	return Copy(e).(*errorx)
}

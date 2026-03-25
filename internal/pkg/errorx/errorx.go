package errorx

import (
	"errors"
	"fmt"

	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

// Type for error handling
type errorx struct {
	err            error
	message        *string
	exception      exceptionx.Exception
	httpStatusCode *int
	parent         errorxInternal
}

// newByType creates an errorx object based on the given type.
func newByType(typ errorType) errorxInternal {
	switch typ {
	case errorTypeWrap:
		errx := &wrapErrorx{}
		errx.parent = errx
		return errx
	case errorTypeJoin:
		errx := &joinErrorx{}
		errx.parent = errx
		return errx
	default:
		return &errorx{}
	}
}

// newErrx creates an errorx object with the given text and type.
func newErrx(text string, typ errorType) errorxInternal {
	errx := newByType(typ)
	errx.setErr(errors.New(text))
	return errx
}

// newFromErr creating errorx from error
func newFromErr(err error) errorxInternal {
	errx := newByType(typeOf(err))

	if err == nil {
		return errx
	}

	if e, ok := err.(errorxInternal); ok {
		errx.initFromErrorx(e)
	} else {
		errx.setErr(err)
	}

	if errx.Exception() == nil {
		errs := UnwrapAll(err)
		if len(errs) > 1 {
			combErrx := WrapErrors(errs...)
			if e, ok := combErrx.(Errorx); ok {
				errx.setException(e.Exception())
			}
		}
	}

	return errx.normalizeToInnerType()
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
func New(text string, args ...any) Errorx {
	var (
		errx    errorxInternal
		errType = errorTypeMain
	)

	for _, arg := range args {
		if _, ok := arg.(error); ok {
			errType = errorTypeWrap
			break
		}
	}

	errx = newErrx(text, errType)
	errx.initFromArgs(args...)

	return errx.normalizeToInnerType()
}

// NewHTTP creates a new Errorx instance with the provided message and HTTP status code.
// Additional optional arguments can be passed to further configure the error.
// The httpStatusCode is automatically added to the arguments as HttpStatusCodeArg.
// Supported argument types and their effects:
//   - exceptionx.Exception: sets the underlying exception
//   - ErrorMessageProvider: overrides the default message
func NewHTTP(text string, httpStatusCode int, args ...any) Errorx {
	var (
		errx    errorxInternal
		errType = errorTypeMain
	)

	for _, arg := range args {
		_, ok := arg.(error)
		if ok {
			errType = errorTypeWrap
			break
		}
	}

	errx = newErrx(text, errType)

	args = append(args, HttpStatusArg(httpStatusCode))
	args = append(args, ErrorMessageArg(text))
	errx.initFromArgs(args...)

	return errx.normalizeToInnerType()
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
	errx.initFromArgs(args...)

	return errx.normalizeToInnerType()
}

// Errorf creates a new Errorx error with a formatted message, similar to fmt.Errorf,
// while also supporting special errorx-specific arguments such as DomainException,
// Exception, ErrorMessageProvider, and HttpStatusProvider.
func Errorf(format string, args ...any) Errorx {
	// Separate arguments into two categories:
	// - special arguments handled by errorx package
	// - regular arguments passed to fmt.Errorf for formatting
	var (
		specialArgs []any // special errorx arguments
		fmtArgs     []any // arguments for standard formatting
	)

	// Pre-allocate fmtArgs only if there are any arguments to avoid unnecessary allocation
	if len(args) > 0 {
		fmtArgs = make([]any, 0, len(args))
	}

	// Partition the input arguments
	for _, arg := range args {
		if isErrorxSpecialArg(arg) {
			specialArgs = append(specialArgs, arg)
		} else {
			fmtArgs = append(fmtArgs, arg)
		}
	}

	// Create the base error using the standard library formatter
	baseErr := fmt.Errorf(format, fmtArgs...)

	// Create a new Errorx instance wrapping the base error
	errx := newByType(typeOf(baseErr))
	errx.setErr(baseErr)

	// Apply any special errorx arguments (e.g. exceptions, status codes, custom messages)
	if len(specialArgs) > 0 {
		errx.initFromArgs(specialArgs...)
	}

	return errx
}

func isErrorxSpecialArg(arg any) bool {
	if arg == nil {
		return false
	}

	switch arg.(type) {
	case exceptionx.DomainException,
		exceptionx.Exception,
		ErrorMessageProvider,
		HttpStatusProvider:
		return true
	}

	return false
}

// parentOrSelf returns an Errorx instance, giving priority to the explicitly set parent.
// If no parent is set, it returns the current errorx object (which implements the Errorx interface).
func (e *errorx) parentOrSelf() errorxInternal {
	if e.parent != nil {
		return e.parent
	}
	return e
}

// initFromErrorx initialization of fields from errorx
func (e *errorx) initFromErrorx(err errorxInternal) {
	e.err = err.Err()
	e.initFromArgs(err.args()...)
}

// initFromArgs initialize errorx from arguments
// args can have types: Exception, ErrorMessageProvider, error
func (e *errorx) initFromArgs(args ...any) {
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
			if err == nil {
				err = v
			} else {
				err = WrapErrors(err, v)
			}
		}
	}

	if message != nil {
		e.message = message
	}

	if exception != nil {
		e.exception = exception
	}

	if httpStatusCode != nil && *httpStatusCode != 0 {
		e.httpStatusCode = httpStatusCode
	}

	if err != nil {
		e.Wrap(err)
	}
}

func (e *errorx) args() []any {
	args := make([]any, 0, 3)

	if e.Exception() != nil {
		args = append(args, e.Exception())
	}
	if text := e.Message(); text != nil {
		args = append(args, ErrorMessageArg(*text))
	}
	if code := e.httpStatusCode; code != nil {
		args = append(args, HttpStatusArg(*code))
	}

	return args
}

// normalizeToInnerType makes the error match the concrete type of its wrapped error.
// Returns the original error if types already match.
func (e *errorx) normalizeToInnerType() errorxInternal {
	if typeOf(e.parentOrSelf()) == typeOf(e.err) {
		return e.parentOrSelf()
	}

	newErr := newByType(typeOf(e.err))
	newErr.initFromErrorx(e.parentOrSelf())

	return newErr
}

// Error returns the error text
func (e *errorx) Error() string {
	if e.err != nil {
		return e.err.Error()
	}
	return ""
}

// Err returns a value of type error
func (e *errorx) Err() error {
	return e.err
}

// setErr sets a value of type error
func (e *errorx) setErr(err error) {
	e.err = err
}

// Message returns a value of type string
func (e *errorx) Message() *string {
	return e.message
}

// Exception returns a exception
func (e *errorx) Exception() exceptionx.Exception {
	return e.exception
}

// setException sets a exception
func (e *errorx) setException(exception exceptionx.Exception) {
	e.exception = exception
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

	exception := UnwrapException(e.parentOrSelf())

	// if exception := UnwrapException(e); exception != nil {
	if exception != nil {
		code := exception.HttpStatusCode()
		if code != 0 {
			return &code
		}
	}

	return nil
}

// Append adds the given errors to the current error using the library's
// wrapping mechanism.
//
// The receiver is mutated in place: only the internal error field is updated.
// The returned Errorx value is the same object as the original (the pointer
// remains unchanged). Nil errors are ignored.
func (e *errorx) Wrap(errs ...error) Errorx {
	e.err = e.WrapAndMerge(errs...).Err()
	return e.normalizeToInnerType()
}

// WrapAndMerge returns a new error that combines the current error with the given
// ones using the library's wrapping mechanism.
//
// A new Errorx object is always returned; the original object remains
// unchanged. Nil errors are ignored.
func (e *errorx) WrapAndMerge(errs ...error) Errorx {
	return WrapErrors(append([]error{e.parentOrSelf()}, errs...)...).(Errorx)
}

// Join appends the given errors to the current error using errors.Join from
// the standard library.
//
// The receiver is mutated in place: the internal error field is replaced
// with the result of errors.Join. The returned Errorx value is the same
// object as the original (the pointer remains unchanged).
// Nil errors in the argument list are ignored.
func (e *errorx) Join(errs ...error) Errorx {
	all := append([]error{e.err}, errs...)
	e.err = errors.Join(all...)
	return e.normalizeToInnerType()
}

// UnwrapAll error analysis errors in a slice
// Duplicates are excluded
func (e *errorx) UnwrapAll() []error {
	return UnwrapAll(e.parentOrSelf())
}

// UnwrapTexts analyzes errors and then extracts text one by one
func (e *errorx) UnwrapTexts() *ErrorTexts {
	return NewErrorTexts().AddUnwrapErr(e.err)
}

// Copy copying an error including nested
func (e *errorx) Copy() error {
	return Copy(e.parentOrSelf())
}

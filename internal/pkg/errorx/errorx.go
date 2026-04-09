package errorx

import (
	"errors"
	"fmt"

	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
	"github.com/neosy/elengrab/internal/pkg/stringx"
)

// errorx is the core struct that implements the Errorx interface.
type errorx struct {
	// err is the underlying error that this errorx wraps. The value is not zero.
	// it may contain a wrapped error.
	err error
	// message is an optional custom message that can be set for the error.
	// If not set, it may be derived from the underlying error or exception.
	// The error text takes precedence over err.
	message *string
	// exception is an optional exception associated with this error.
	// It can provide additional context or metadata about the error.
	exception exceptionx.Exception
	// httpStatusCode is an optional HTTP status code associated with this error.
	// Takes precedence over the code in Exception.
	httpStatusCode *int
	// parent is used by embedding structs and holds a reference to the outer struct.
	// It is required to determine the actual concrete type when methods are invoked
	// from the embedded errorx.
	//
	// For example:
	//
	//	type wrapErrorx struct {
	//	    errorx
	//	}
	//
	// For this structure, parent will have the type wrapErrorx.
	// If errorx is used directly (without embedding), parent will be of type errorx.
	parent errorxInternal
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

// NewFromDomainException creates a new Errorx instance based on the provided DomainException.
func NewFromDomainException(exception exceptionx.DomainException, args ...any) Errorx {
	return New(exception.Message(), exception)
}

// isErrorxSpecialArg determines if an argument is a special type
// that should be processed by errorx's argument handling logic.
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
		msg := stringx.Capitalize(*message)
		e.message = &msg
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

// args returns the arguments used to configure the error.
func (e *errorx) args() []any {
	args := make([]any, 0, 3)

	if e.Exception() != nil {
		args = append(args, e.Exception())
	}
	if text := e.Message(); text != "" {
		args = append(args, ErrorMessageArg(text))
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
func (e *errorx) Message() string {
	if e.message == nil {
		return ""
	}
	return *e.message
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
func (e *errorx) HttpStatusCodeRaw() int {
	if e.httpStatusCode == nil {
		return 0
	}
	return *e.httpStatusCode
}

// HttpStatusCode returns the HTTP status code, falling back to the exception if needed.
func (e *errorx) HttpStatusCode() int {
	if e.httpStatusCode != nil {
		return *e.httpStatusCode
	}

	exception := OuterException(e.parentOrSelf())

	if exception != nil {
		code := exception.HttpStatusCode()
		if code != 0 {
			return code
		}
	}

	return 0
}

// Wrap adds the given errors to the current error in place using the library's wrapping mechanism.
// The current Errorx object is mutated; the receiver remains the same. Nil errors are ignored.
func (e *errorx) Wrap(errs ...error) Errorx {
	e.err = e.WrapNew(errs...).Err()
	return e.normalizeToInnerType()
}

// WrapNew returns a new Errorx that wraps the current error along with the provided errors.
// The original error remains unchanged; nil errors are ignored.
func (e *errorx) WrapNew(errs ...error) Errorx {
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
// Duplicates are not excluded
func (e *errorx) UnwrapAll() []error {
	return UnwrapAll(e.parentOrSelf())
}

// UnwrapUnique error analysis errors in a slice
// Duplicates are excluded
func (e *errorx) UnwrapUnique() []error {
	return UnwrapUnique(e.parentOrSelf())
}

// UnwrapTexts analyzes errors and then extracts text one by one
func (e *errorx) UnwrapTexts() *ErrorTexts {
	return NewErrorTexts().AddUnwrapErr(e.err)
}

// Copy copying an error including nested
func (e *errorx) Copy() error {
	return Copy(e.parentOrSelf())
}

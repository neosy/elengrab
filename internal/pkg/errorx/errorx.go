package errorx

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
	"github.com/neosy/elengrab/internal/pkg/stringx"
	"github.com/valyala/fasthttp"
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
	// httpStatus is an optional HTTP status associated with this error.
	// Takes precedence over the status in Exception.
	httpStatus *int
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
//   - HttpStatusProvider: extracts and sets the HTTP status
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

// NewWithMessage creates a new Errorx instance with the provided text and an explicit message.
// The message is set using the WithErrorMessage argument, which takes precedence over the default message.
// Additional optional arguments can be passed to further configure the error.
// Supported argument types and their effects:
//   - exceptionx.Exception: sets the underlying exception
//   - HttpStatusProvider: extracts and sets the HTTP status
//   - error: will be combined with the current error using CombineErrors
//
// The provided text is used to create the base error, while the message argument allows for a custom message to be set explicitly.
func NewWithMessage(text string, args ...any) Errorx {
	return New(text, append([]any{WithErrorMessage(text)}, args...)...)
}

// NewHTTP creates a new Errorx instance with the provided text error and HTTP status.
// Additional optional arguments can be passed to further configure the error.
// The httpStatus is automatically added to the arguments as httpStatusArg.
// Supported argument types and their effects:
//   - exceptionx.Exception: sets the underlying exception
//   - ErrorMessageProvider: overrides the default message
func NewHTTP(text string, httpStatus int, args ...any) Errorx {
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

	if text == "" {
		text = stringx.LowerFirst(http.StatusText(httpStatus))
	}

	errx = newErrx(text, errType)
	errx.initFromArgs(append(args, WithHttpStatus(httpStatus))...)

	return errx.normalizeToInnerType()
}

// NewHTTPStatus creates a new Errorx instance with the provided HTTP status and an optional text error.
// If the text is empty, it defaults to the standard HTTP status text corresponding to the provided status code.
// Additional optional arguments can be passed to further configure the error.
// Supported argument types and their effects:
//   - exceptionx.Exception: sets the underlying exception
//   - ErrorMessageProvider: overrides the default message
func NewHTTPStatus(status int, args ...any) Errorx {
	return NewHTTP("", status, args...)
}

// NewHTTPMessage creates a new Errorx instance with the provided text error and an HTTP status derived from the text.
// If the text is empty, it defaults to a 500 Internal Server Error status.
// Additional optional arguments can be passed to further configure the error.
// Supported argument types and their effects:
//   - exceptionx.Exception: sets the underlying exception
//   - HttpStatusProvider: extracts and sets the HTTP status (overrides default status derived from text)
func NewHTTPMessage(text string, httpStatus int, args ...any) Errorx {
	return NewHTTP(text, httpStatus, append([]any{WithErrorMessage(text)}, args...)...)
}

// NewFromError creating errorx from error and arguments
// args can have types: Exception, ErrorMessageProvider, error
// Supported argument types and their effects:
//   - exceptionx.Exception: sets the underlying exception
//   - HttpStatusProvider: extracts and sets the HTTP status
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

	// Apply any special errorx arguments (e.g. exceptions, status statuss, custom messages)
	if len(specialArgs) > 0 {
		errx.initFromArgs(specialArgs...)
	}

	return errx
}

// NewFromException creates a new Errorx instance based on the provided exception.
func NewFromException(ex exceptionx.Exception, args ...any) Errorx {
	return New(ex.Error(), append([]any{ex}, args...)...)
}

// NewFromDomainException creates a new Errorx instance based on the provided DomainException.
func NewFromDomainException(ex exceptionx.DomainException, args ...any) Errorx {
	return New(stringx.LowerFirst(ex.Message()), append([]any{ex}, args...)...)
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
	var httpStatus *int
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
				httpStatus = v()
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

	if httpStatus != nil && *httpStatus != 0 {
		e.httpStatus = httpStatus
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
		args = append(args, WithErrorMessage(text))
	}
	if status := e.httpStatus; status != nil {
		args = append(args, WithHttpStatus(*status))
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

// PublicMessage returns the public-facing error message, prioritizing the explicitly set message,
// then the message from the associated exception,
// and finally falling back to the standard HTTP status text if available.
func (e *errorx) PublicMessage() string {
	msg := e.OuterMessage()
	if msg != "" {
		return stringx.Capitalize(msg)
	}

	exception := e.OuterException()
	if exception != nil {
		if msg := exception.Message(); msg != "" {
			return stringx.Capitalize(msg)
		}
	}

	return stringx.Capitalize(e.PublicHttpStatusText())
}

// HttpStatus returns the explicitly set HTTP status.
func (e *errorx) HttpStatus() int {
	if e.httpStatus == nil {
		return 0
	}
	return *e.httpStatus
}

// PublicHttpStatus returns the HTTP status, prioritizing the explicitly set status,
// then the status from the associated exception, and finally falling back to 500 if none is set.
func (e *errorx) PublicHttpStatus() int {
	if e.httpStatus != nil {
		return *e.httpStatus
	}

	// Check for HTTP status in the error chain, giving priority to the outermost error.
	err := outerErrorxWithExceptionOrHttpStatus(e.parentOrSelf())
	if err != nil {
		if status := err.HttpStatus(); status != 0 {
			return status
		}
		if err.Exception() != nil {
			if status := err.Exception().HttpStatus(); status != 0 {
				return status
			}
		}
	}

	return fasthttp.StatusInternalServerError
}

// PublicHttpStatusText returns the standard HTTP status text corresponding to the public HTTP status.
func (e *errorx) PublicHttpStatusText() string {
	return http.StatusText(e.PublicHttpStatus())
}

// Exception returns a exception
func (e *errorx) Exception() exceptionx.Exception {
	return e.exception
}

// PublicException returns the first (outermost) exception found in the error chain.
// It is a public-facing method that provides access to the exception associated with the error, if any.
func (e *errorx) PublicException() exceptionx.Exception {
	return e.OuterException()
}

// setException sets a exception
func (e *errorx) setException(exception exceptionx.Exception) {
	e.exception = exception
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

// OuterMessage returns the first (outermost) message found in the error chain.
func (e *errorx) OuterMessage() string {
	return OuterMessage(e.parentOrSelf())
}

// RootMessage returns the first (outermost) message found in the error chain.
func (e *errorx) RootMessage() string {
	return RootMessage(e.parentOrSelf())
}

// OuterException returns the first (outermost) exception found in the error chain.
func (e *errorx) OuterException() exceptionx.Exception {
	return OuterException(e.parentOrSelf())
}

// RootException returns the first (outermost) exception found in the error chain.
func (e *errorx) RootException() exceptionx.Exception {
	return RootException(e.parentOrSelf())
}

// OuterHttpStatus returns the first (outermost) HTTP status found in the error chain.
func (e *errorx) OuterHttpStatus() int {
	return OuterHttpStatus(e.parentOrSelf())
}

// RootHttpStatus returns the first (outermost) HTTP status found in the error chain.
func (e *errorx) RootHttpStatus() int {
	return RootHttpStatus(e.parentOrSelf())
}

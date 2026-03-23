package errorx

import (
	"errors"

	"github.com/neosy/elengrab/pkg/errorx/exceptionx"
)

// Type for error handling
type Errorx struct {
	err            error
	message        ErrorxMessage
	exception      exceptionx.Exception
	httpStatusCode *int
}

// newErrx creating an Errorx object from text
func newErrx(text string) (errx ErrorxInterface) {
	errx = &Errorx{
		err: errors.New(text),
	}

	return
}

// newFromErr creating Errorx from error
func newFromErr(err error) ErrorxInterface {
	errx := new(Errorx)

	if err == nil {
		return errx
	}

	if e, ok := err.(ErrorxInterface); ok {
		errx.initFromErrorx(e)
	} else {
		errx.err = err
	}

	return errx
}

// New creating an Errorx object from text and arguments
// args can have types: ExceptionType, ExceptionCode, ErrorxMessage, error, uint
//
//	message: ErrorxMessage
//	new num ExceptionCode: uint
func New(text string, args ...any) (errx ErrorxInterface) {
	errx = newErrx(text)

	errx.(*Errorx).initFromArgs(args...)

	return
}

// NewByErr creating Errorx from error and arguments
// args can have types: ExceptionType, ExceptionCode, ErrorxMessage, error, uint
//
//	message: ErrorxMessage
//	new num ExceptionCode: uint
func NewByErr(err error, args ...any) ErrorxInterface {
	if err == nil {
		return nil
	}

	errx := newFromErr(err)

	errx.(*Errorx).initFromArgs(args...)

	return errx
}

// NewByExceptionType creating Errorx from text and ExceptionType
func NewByExceptionType(text string, eType exceptionx.ExceptionType) (errx ErrorxInterface) {
	errx = New(text, eType)

	return
}

// NewByExceptionCode creating Errorx from ExceptionCode
func NewByExceptionCode(code exceptionx.ExceptionCode) (errx ErrorxInterface) {
	errx = New(code.String(), code)

	return
}

// NewDomainException creating Errorx from text and ExceptionType
func NewDomainException(text string, eType exceptionx.ExceptionType, num uint) (errx ErrorxInterface) {
	code := exceptionx.NewExceptionCode(num, text, eType)

	errx = NewByExceptionCode(code)

	return
}

// initFromErrorx initialization of fields from Errorx
func (errx *Errorx) initFromErrorx(err ErrorxInterface) {
	errx.err = err.Err()

	args := make([]any, 4)
	args = append(args, err.Message(), err.ExceptionCode(), err.ExceptionType())

	var httStatusCode *int
	e, ok := err.(*Errorx)
	if ok {
		httStatusCode = e.httpStatusCode
	} else {
		httStatusCode = err.HttpStatusCode()
	}

	if httStatusCode != nil {
		args = append(args, ArgHttpStatusCode(*httStatusCode))
	}

	errx.initFromArgs(args...)
}

// initFromArgs initialize Errorx from arguments
// args can have types: ExceptionType, ExceptionCode, string, ErrorxMessage, error
func (errx *Errorx) initFromArgs(args ...any) {
	var message ErrorxMessage
	var eType exceptionx.ExceptionType
	var eCode exceptionx.ExceptionCode
	var httpStatusCode *int
	var err error

	for _, arg := range args {
		switch v := arg.(type) {
		case exceptionx.ExceptionType:
			eType = v
		case exceptionx.ExceptionCode:
			eCode = v
		case HttpStatusProvider:
			code := v()
			httpStatusCode = &code
		case ErrorxMessage:
			message = v
		case error:
			err = CombineErrors(err, v)
		}
	}

	if errx.ExceptionCode() == nil {
		var eCodeNum *uint
		for _, arg := range args {
			switch v := arg.(type) {
			case uint:
				eCodeNum = &v
			}
		}

		if eCodeNum != nil {
			eCode = exceptionx.NewExceptionCode(*eCodeNum, errx.Error(), errx.ExceptionType())
		}
	}

	if message != nil {
		errx.message = message
	}

	if eType == nil && errx.ExceptionType() != nil {
		eType = errx.ExceptionType()
	}

	if eCode == nil && errx.ExceptionCode() != nil {
		eCode = errx.ExceptionCode()
	}

	errx.exception = *exceptionx.NewException(eCode, eType)
	errx.httpStatusCode = httpStatusCode

	if err != nil {
		errx.Append(err)
	}
}

// Error returns the error text
func (e *Errorx) Error() (text string) {
	if e.err != nil {
		text = e.err.Error()
	}

	return
}

// Err returns a value of type error
func (e *Errorx) Err() error {
	return e.err
}

// Message returns a value of type ErrorxMessage
func (e *Errorx) Message() ErrorxMessage {
	return e.message
}

// ExceptionType returns a type of exception
func (e *Errorx) ExceptionType() exceptionx.ExceptionType {
	return e.exception.Type()
}

// ExceptionCode returns a code of exception
func (e *Errorx) ExceptionCode() exceptionx.ExceptionCode {
	return e.exception.Code()
}

// HttpStatusCode returns the HTTP status code associated with the error, if any.
func (e *Errorx) HttpStatusCode() *int {
	if e.httpStatusCode != nil {
		return e.httpStatusCode
	}

	if e.exception.Type() != nil {
		code := e.exception.Type().HttpStatusCode()
		if code != 0 {
			return &code
		}
	}

	if e.exception.Code() != nil && e.exception.Code().Type() != nil {
		code := e.exception.Code().Type().HttpStatusCode()
		if code != 0 {
			return &code
		}
	}

	return nil
}

// Append merges the given errors into the current error.
// The pointer remains the same; only the internal field err are updated.
func (e *Errorx) Append(errs ...error) ErrorxInterface {
	e.err = e.Combine(errs...).Err()
	return e
}

// Combine returns a new error that combines the current error with the given ones.
// A new object is returned; the current object remains unchanged.
func (e *Errorx) Combine(errs ...error) ErrorxInterface {
	return CombineErrors(e, CombineErrors(errs...)).(ErrorxInterface)
}

// UnwrapAll error analysis errors in a slice
// Duplicates are excluded
func (e *Errorx) UnwrapAll() []error {
	return UnwrapAll(e)
}

// UnwrapTexts analyzes errors and then extracts text one by one
func (e *Errorx) UnwrapTexts() *ErrorTexts {
	return NewErrorTexts().AddUnwrapErr(e)
}

// Copy copying an error including nested
func (e *Errorx) Copy() ErrorxInterface {
	return Copy(e).(*Errorx)
}

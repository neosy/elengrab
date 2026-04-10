package exceptionx

import (
	"github.com/neosy/elengrab/internal/pkg/errorx/internal/types"
	"github.com/neosy/elengrab/internal/pkg/stringx"
	"github.com/valyala/fasthttp"
)

type Exception interface {
	// Num returns the exception number.
	Num() uint
	// Code returns the exception code (short text).
	Code() string
	// String returns the exception text. If the message is empty, it returns the code.
	String() string
	// Message returns the exception message (long text).
	Message() string
	// Error implements the error interface by returning the exception message.
	Error() string
	// HttpStatus returns the HTTP status associated with the exception.
	HttpStatus() int
}

// Type Exception
type exception struct {
	num        uint
	code       string
	message    string
	httpStatus int
}

// NewException creating an Exception object
func NewException(num uint, args ...any) Exception {
	newArgs := ExceptionArgs{}

	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			newArgs.Message = v
		case ExceptionArg:
			if v != nil {
				v(&newArgs)
			}
		case types.HttpStatusProvider:
			if v != nil {
				if status := v(); status != nil {
					newArgs.HTTPStatus = *status
				}
			}
		case Exception:
			newArgs.Code = v.Code()
			newArgs.Message = v.Message()
			newArgs.HTTPStatus = v.HttpStatus()
		}
	}

	if newArgs.HTTPStatus == 0 {
		newArgs.HTTPStatus = fasthttp.StatusInternalServerError
	}

	ex := &exception{
		num:        num,
		code:       newArgs.Code,
		message:    newArgs.Message,
		httpStatus: newArgs.HTTPStatus,
	}

	return ex
}

// Num returns the exception number.
func (e *exception) Num() uint {
	return e.num
}

// Code returns the exception code.
func (e *exception) Code() string {
	if e == nil {
		return ""
	}
	return e.code
}

// String returns the exception text. If the message is empty, it returns the code.
func (e *exception) String() string {
	if e == nil {
		return ""
	}
	text := stringx.Capitalize(e.message)
	if text == "" {
		text = e.code
	}
	return text
}

// String returns the exception text.
func (e *exception) Message() string {
	if e == nil {
		return ""
	}
	return stringx.Capitalize(e.message)
}

// Error implements the error interface by returning the exception message.
func (e *exception) Error() string {
	if e == nil {
		return ""
	}
	return stringx.LowerFirst(e.message)
}

// HttpStatus returns the HTTP status associated with the exception.
func (e *exception) HttpStatus() int {
	return e.httpStatus
}

// NewErrorx creates a new Errorx instance based on this exception and the provided arguments.
func (e *exception) NewErrorx(args ...any) error {
	return errorxFactory.NewFromException(e, args...)
}

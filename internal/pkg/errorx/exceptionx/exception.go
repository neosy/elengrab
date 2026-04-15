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
	// HTTPStatus returns the HTTP status associated with the exception.
	HTTPStatus() int
	// NewErrorx creates a new Errorx instance based on this exception and the provided arguments.
	NewErrorx(args ...any) error
}

// Type Exception
type exception struct {
	num        uint
	code       string
	message    string
	httpStatus int
}

// NewException creating an Exception object
func NewException(num uint, code string, opts ...any) Exception {
	newOpts := ExceptionOptions{}

	for _, opt := range opts {
		switch v := opt.(type) {
		case string:
			newOpts.Message = v
		case ExceptionOption:
			if v != nil {
				v(&newOpts)
			}
		case types.HttpStatusProvider:
			if v != nil {
				if status := v(); status != nil {
					newOpts.HTTPStatus = *status
				}
			}
		case Exception:
			newOpts.Message = v.Message()
			newOpts.HTTPStatus = v.HTTPStatus()
		}
	}

	if newOpts.HTTPStatus == 0 {
		newOpts.HTTPStatus = fasthttp.StatusInternalServerError
	}

	ex := &exception{
		num:        num,
		code:       code,
		message:    newOpts.Message,
		httpStatus: newOpts.HTTPStatus,
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
	text := e.Message()
	if text != "" {
		return text
	}
	return e.code
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

// HTTPStatus returns the HTTP status associated with the exception.
func (e *exception) HTTPStatus() int {
	return e.httpStatus
}

// NewErrorx creates a new Errorx instance based on this exception and the provided arguments.
func (e *exception) NewErrorx(args ...any) error {
	return errorxFactory.NewFromException(e, args...)
}

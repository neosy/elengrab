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
	// HttpStatusCode returns the HTTP status code associated with the exception.
	HttpStatusCode() int
}

// Type Exception
type exception struct {
	num            uint
	code           string
	message        string
	httpStatusCode int
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
				if code := v(); code != nil {
					newArgs.HTTPStatusCode = *code
				}
			}
		case Exception:
			newArgs.Code = v.Code()
			newArgs.Message = v.Message()
			newArgs.HTTPStatusCode = v.HttpStatusCode()
		}
	}

	if newArgs.HTTPStatusCode == 0 {
		newArgs.HTTPStatusCode = fasthttp.StatusInternalServerError
	}

	ex := &exception{
		num:            num,
		code:           newArgs.Code,
		message:        newArgs.Message,
		httpStatusCode: newArgs.HTTPStatusCode,
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

// HttpStatusCode returns the HTTP status code associated with the exception.
func (e *exception) HttpStatusCode() int {
	return e.httpStatusCode
}

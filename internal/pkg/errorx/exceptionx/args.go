package exceptionx

import "github.com/neosy/elengrab/internal/pkg/errorx/internal/types"

type (
	HttpStatusProvider = types.HttpStatusProvider

	ExceptionArgs struct {
		Code       string
		Message    string
		HTTPStatus int
	}
	ExceptionArg func(args *ExceptionArgs)
)

var (
	WithHttpStatus   = types.WithHttpStatus
	WithErrorMessage = types.WithErrorMessage
)

// ExceptionArgCode creates an argument that sets the exception code.
func ExceptionArgCode(code string) ExceptionArg {
	return func(args *ExceptionArgs) {
		args.Code = code
	}
}

// ExceptionArgMessage creates an argument that sets the exception message.
func ExceptionArgMessage(message string) ExceptionArg {
	return func(args *ExceptionArgs) {
		args.Message = message
	}
}

// ExceptionArgHTTPStatus creates an argument that sets the HTTP status for the exception.
func ExceptionArgHTTPStatus(httpStatus int) ExceptionArg {
	return func(args *ExceptionArgs) {
		args.HTTPStatus = httpStatus
	}
}

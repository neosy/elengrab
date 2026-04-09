package exceptionx

import "github.com/neosy/elengrab/internal/pkg/errorx/internal/types"

type (
	HttpStatusProvider = types.HttpStatusProvider

	ExceptionArgs struct {
		Code           string
		Message        string
		HTTPStatusCode int
	}
	ExceptionArg func(args *ExceptionArgs)
)

var (
	HttpStatusArg = types.HttpStatusArg
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

// ExceptionArgHTTPStatusCode creates an argument that sets the HTTP status code for the exception.
func ExceptionArgHTTPStatusCode(httpStatusCode int) ExceptionArg {
	return func(args *ExceptionArgs) {
		args.HTTPStatusCode = httpStatusCode
	}
}

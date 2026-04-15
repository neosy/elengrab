package exceptionx

import "github.com/neosy/elengrab/internal/pkg/errorx/internal/types"

type (
	HttpStatusProvider = types.HttpStatusProvider

	ExceptionOptions struct {
		Message    string
		HTTPStatus int
	}
	ExceptionOption func(opts *ExceptionOptions)
)

// WithMessage creates an argument that sets the exception message.
func WithMessage(message string) ExceptionOption {
	return func(opts *ExceptionOptions) {
		opts.Message = message
	}
}

// WithHTTPStatus creates an argument that sets the HTTP status for the exception.
func WithHTTPStatus(httpStatus int) ExceptionOption {
	return func(opts *ExceptionOptions) {
		opts.HTTPStatus = httpStatus
	}
}

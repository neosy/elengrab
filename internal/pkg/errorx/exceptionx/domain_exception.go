package exceptionx

import (
	"net/http"

	"github.com/neosy/elengrab/internal/pkg/errorx/internal/types"
	"github.com/valyala/fasthttp"
)

// DomainException represents a domain-specific exception with associated HTTP status and error message.
// It provides methods to retrieve the exception number, text, message, error message provider, HTTP status,
// and to create a new Exception based on the domain exception type.
type DomainException interface {
	// Num exception number
	Num() uint
	// Code returns the exception code, which is the same as the string representation of the exception.
	Code() string
	// Message exception message
	Message() string
	// String returns the exception code (short text).
	String() string
	// ErrorMessage returns error message provider
	ErrorMessage() types.ErrorMessageProvider
	// HTTPStatus returns HTTP status
	HTTPStatus() int
	// NewException creates a new exception with the given domain exception type.
	NewException() Exception
	// NewErrorx creates a new Errorx instance based on this domain exception and the provided arguments.
	NewErrorx(args ...any) error
}

// domainException type of basic exception
type domainException uint

const (
	ERROR domainException = iota + 101
	WRONG_DATA
	VALIDATE
	UNAUTHORIZED
	NOT_FOUND
	METHOD_NOT_ALLOWED
	TOO_MANY_REQUESTS
	MAINTENANCE
	FORBIDDEN
	NO_CONTENT
)

// Map of exception types and their short code
var domainExceptionCodeMap = map[domainException]string{
	ERROR:              "ERROR",
	WRONG_DATA:         "WRONG_DATA",
	VALIDATE:           "VALIDATE",
	UNAUTHORIZED:       "UNAUTHORIZED",
	NOT_FOUND:          "NOT_FOUND",
	METHOD_NOT_ALLOWED: "METHOD_NOT_ALLOWED",
	TOO_MANY_REQUESTS:  "TOO_MANY_REQUESTS",
	MAINTENANCE:        "MAINTENANCE",
	FORBIDDEN:          "FORBIDDEN",
	NO_CONTENT:         "NO_CONTENT",
}

var domainExceptionMessageMap = map[domainException]string{
	WRONG_DATA: "Bad Request", // 400
	VALIDATE:   "Bad Request", // 400
}

// Map of exception types and HTTP statuss
var domainExceptionHttpStatusMap = map[domainException]int{
	ERROR:              fasthttp.StatusInternalServerError,
	WRONG_DATA:         fasthttp.StatusBadRequest,
	VALIDATE:           fasthttp.StatusBadRequest,
	UNAUTHORIZED:       fasthttp.StatusUnauthorized,
	NOT_FOUND:          fasthttp.StatusNotFound,
	METHOD_NOT_ALLOWED: fasthttp.StatusMethodNotAllowed,
	TOO_MANY_REQUESTS:  fasthttp.StatusTooManyRequests,
	MAINTENANCE:        fasthttp.StatusServiceUnavailable,
	FORBIDDEN:          fasthttp.StatusForbidden,
	NO_CONTENT:         fasthttp.StatusNoContent,
}

// Num exception number
func (e domainException) Num() uint {
	return uint(e)
}

// Code returns the exception code, which is the same as the string representation of the exception.
func (e domainException) Code() string {
	return domainExceptionCodeMap[e]
}

// String returns the exception code (short text).
func (e domainException) String() string {
	return e.Code()
}

// Message exception message
func (e domainException) Message() string {
	msg := domainExceptionMessageMap[e]
	if msg == "" {
		msg = http.StatusText(e.HTTPStatus())
	}
	return msg
}

// ErrorMessage returns error message provider
func (e domainException) ErrorMessage() types.ErrorMessageProvider {
	return types.WithErrorMessage(e.Message())
}

// HTTPStatus returns HTTP status
func (e domainException) HTTPStatus() int {
	status, ok := domainExceptionHttpStatusMap[e]
	if !ok {
		return fasthttp.StatusInternalServerError
	}
	return status
}

// NewException creates a new exception with the given domain exception type.
func (e domainException) NewException() Exception {
	return NewException(
		e.Num(),
		e.Code(),
		WithMessage(e.Message()),
		WithHTTPStatus(e.HTTPStatus()),
	)
}

// NewErrorx creates a new Errorx instance based on this domain exception and the provided arguments.
func (e domainException) NewErrorx(args ...any) error {
	return errorxFactory.NewFromDomainException(e, args...)
}

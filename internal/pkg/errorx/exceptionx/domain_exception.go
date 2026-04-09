package exceptionx

import (
	"github.com/neosy/elengrab/internal/pkg/errorx/internal/types"
	"github.com/valyala/fasthttp"
)

// DomainException represents a domain-specific exception with associated HTTP status code and error message.
// It provides methods to retrieve the exception number, text, message, error message provider, HTTP status code,
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
	// HttpStatusCode returns HTTP status code
	HttpStatusCode() int
	// NewException creates a new exception with the given domain exception type.
	NewException() Exception
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

// Map of exception types and their text descriptions
var domainExceptionTextMap = map[domainException]string{
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
	ERROR:              "internal Server Error", // 500
	WRONG_DATA:         "bad Request",           // 400
	VALIDATE:           "bad Request",           // 400
	UNAUTHORIZED:       "unauthorized",          // 401
	NOT_FOUND:          "not Found",             // 404
	METHOD_NOT_ALLOWED: "method Not Allowed",    // 405
	TOO_MANY_REQUESTS:  "too Many Requests",     // 429
	MAINTENANCE:        "service Unavailable",   // 503
	FORBIDDEN:          "forbidden",             // 403
	NO_CONTENT:         "no Content",            // 204
}

// Map of exception types and HTTP status codes
var domainExceptionHttpStatusCodeMap = map[domainException]int{
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
func (v domainException) Num() uint {
	return uint(v)
}

// Code returns the exception code, which is the same as the string representation of the exception.
func (v domainException) Code() string {
	return domainExceptionTextMap[v]
}

// String returns the exception code (short text).
func (v domainException) String() string {
	return v.Code()
}

// Message exception message
func (v domainException) Message() string {
	return domainExceptionMessageMap[v]
}

// ErrorMessage returns error message provider
func (v domainException) ErrorMessage() types.ErrorMessageProvider {
	return types.ErrorMessageArg(v.Message())
}

// HttpStatusCode returns HTTP status code
func (e domainException) HttpStatusCode() int {
	code, ok := domainExceptionHttpStatusCodeMap[e]
	if !ok {
		return fasthttp.StatusInternalServerError
	}
	return code
}

// NewException creates a new exception with the given domain exception type.
func (v domainException) NewException() Exception {
	return NewException(
		v.Num(),
		ExceptionArgCode(v.Code()),
		ExceptionArgMessage(v.Message()),
		ExceptionArgHTTPStatusCode(v.HttpStatusCode()),
	)
}

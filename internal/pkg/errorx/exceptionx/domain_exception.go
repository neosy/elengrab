package exceptionx

import (
	"github.com/valyala/fasthttp"
)

// DomainException type of basic exception
type DomainException uint

const (
	ERROR DomainException = iota + 101
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
var domainExceptionTextMap = map[DomainException]string{
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

// Map of exception types and HTTP status codes
var domainExceptionHttpStatusCodeMap = map[DomainException]int{
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
func (v DomainException) Num() uint {
	return uint(v)
}

// String exception text
func (v DomainException) String() string {
	return domainExceptionTextMap[v]
}

// HttpStatusCode returns HTTP status code
func (e DomainException) HttpStatusCode() int {
	code, ok := domainExceptionHttpStatusCodeMap[e]
	if !ok {
		return fasthttp.StatusInternalServerError
	}
	return code
}

// NewException creates a new exception with the given domain exception type.
func (v DomainException) NewException() Exception {
	return NewException(uint(v), v.String(), HttpStatusArg(v.HttpStatusCode()))
}

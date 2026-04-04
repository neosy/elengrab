package exceptionx

import (
	errmsg "github.com/neosy/elengrab/internal/pkg/errorx/internal/error_message"
	"github.com/neosy/elengrab/internal/pkg/errorx/internal/utils"
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

var domainExceptionMessageMap = map[DomainException]string{
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

func (v DomainException) Message() string {
	return domainExceptionMessageMap[v]
}

func (v DomainException) ErrorMessage() errmsg.ErrorMessageProvider {
	return errmsg.ErrorMessageArg(utils.Capitalize(v.Message()))
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

package exceptionx

import (
	"net/http"

	"github.com/valyala/fasthttp"
)

// DomainExceptionType type of basic exception
type DomainExceptionType string

const (
	ERROR              DomainExceptionType = "ERROR"
	WRONG_DATA         DomainExceptionType = "WRONG_DATA"
	VALIDATE           DomainExceptionType = "VALIDATE"
	UNAUTHORIZED       DomainExceptionType = "UNAUTHORIZED"
	NOT_FOUND          DomainExceptionType = "NOT_FOUND"
	BUSINESS           DomainExceptionType = "BUSINESS"
	METHOD_NOT_ALLOWED DomainExceptionType = "METHOD_NOT_ALLOWED"
	RATE_LIMIT         DomainExceptionType = "RATE_LIMIT"
	MAINTENANCE        DomainExceptionType = "MAINTENANCE"
	FORBIDDEN          DomainExceptionType = "FORBIDDEN"
	NO_CONTENT         DomainExceptionType = "NO_CONTENT"
)

// Map of exception types and HTTP status codes
var domainExceptionTypeHttpStatusCodeMap = map[ExceptionType]int{
	ERROR:              fasthttp.StatusInternalServerError,
	WRONG_DATA:         fasthttp.StatusBadRequest,
	VALIDATE:           fasthttp.StatusBadRequest,
	UNAUTHORIZED:       fasthttp.StatusUnauthorized,
	NOT_FOUND:          fasthttp.StatusNotFound,
	BUSINESS:           fasthttp.StatusInternalServerError,
	METHOD_NOT_ALLOWED: fasthttp.StatusMethodNotAllowed,
	RATE_LIMIT:         fasthttp.StatusTooManyRequests,
	MAINTENANCE:        fasthttp.StatusServiceUnavailable,
	FORBIDDEN:          fasthttp.StatusForbidden,
	NO_CONTENT:         fasthttp.StatusNoContent,
}

// String exception text
func (v DomainExceptionType) String() string {
	return string(v)
}

// HttpStatusCode returns HTTP status code
func (e DomainExceptionType) HttpStatusCode() int {
	code, ok := domainExceptionTypeHttpStatusCodeMap[e]
	if !ok {
		code = http.StatusInternalServerError
	}

	return code
}

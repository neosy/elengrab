package exceptionx

import "github.com/valyala/fasthttp"

type Exception interface {
	Num() uint
	String() string
	HttpStatusCode() int
}

// Type Exception
type exception struct {
	num            uint
	text           string
	httpStatusCode int
}

// NewException creating an Exception object
func NewException(num uint, args ...any) Exception {
	ex := &exception{
		num: num,
	}

	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			ex.text = v
		case HttpStatusProvider:
			if v != nil {
				if code := v(); code != nil {
					ex.httpStatusCode = *code
				}
			}
		case Exception:
			ex.text = v.String()
			ex.httpStatusCode = v.HttpStatusCode()
		}
	}

	if ex.httpStatusCode == 0 {
		ex.httpStatusCode = fasthttp.StatusInternalServerError
	}

	return ex
}

// Num returns the exception number.
func (e *exception) Num() uint {
	return e.num
}

// String returns the exception text.
func (e *exception) String() string {
	return e.text
}

// HttpStatusCode returns the HTTP status code associated with the exception.
func (e *exception) HttpStatusCode() int {
	return e.httpStatusCode
}

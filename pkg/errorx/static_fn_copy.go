package errorx

import (
	"errors"
	"fmt"
)

// Copy copying an error including nested
func Copy(err error) (newErr error) {
	if err == nil {
		return
	}

	var unwrap func(err error)

	unwrap = func(err error) {
		switch e := err.(type) {
		case (interface{ Unwrap() []error }):
			for _, subErr := range e.Unwrap() {
				unwrap(subErr)
			}
		case ErrorxInterface:
			args := make([]any, 4)
			args = append(args, e.Message(), e.ExceptionType(), e.ExceptionCode())

			var httpStatusCode *int
			ex, ok := e.(*Errorx)
			if ok {
				httpStatusCode = ex.httpStatusCode
			} else {
				httpStatusCode = e.HttpStatusCode()
			}

			if httpStatusCode != nil {
				args = append(args, ArgHttpStatusCode(*httpStatusCode))
			}

			errx := New(
				e.Error(),
				args...,
			)
			errx.(*Errorx).err = Copy(e.Err())

			if newErr == nil {
				newErr = errx
			} else {
				newErr = fmt.Errorf("%w%s%w", newErr, combineErrorsSeparator(), errx)
			}
		default:
			if newErr == nil {
				newErr = errors.New(e.Error())
			} else {
				newErr = fmt.Errorf("%w%s%w", newErr, combineErrorsSeparator(), e)
			}
		}
	}

	unwrap(err)

	return
}

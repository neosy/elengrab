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
			errx := New(
				e.Error(),
				e.Message(),
				e.ExceptionType(),
				e.ExceptionCode(),
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

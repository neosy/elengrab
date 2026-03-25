package errorx

import (
	"errors"

	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

type wrapErrorx struct {
	errorx
}

func (e *wrapErrorx) Unwrap() error {
	wrapErr, ok := e.Err().(interface{ Unwrap() error })
	if !ok {
		return nil
	}
	return wrapErr.Unwrap()
}

type joinErrorx struct {
	errorx
}

func (e *joinErrorx) Unwrap() []error {
	joinErr, ok := e.Err().(interface{ Unwrap() []error })
	if !ok {
		return nil
	}
	return joinErr.Unwrap()
}

// UnwrapAll error analysis errors in a slice
// Duplicates are excluded
func UnwrapAll(err error) []error {
	errProcessed := make(map[string]struct{})

	var (
		errs   []error
		unwrap func(error)
	)

	unwrap = func(e error) {
		if e == nil {
			return
		}

		errStr := e.Error()
		if _, exists := errProcessed[errStr]; exists {
			return
		}

		errProcessed[errStr] = struct{}{}

		// Checking for multiple nested errors (for errors.Join)
		unwrappedErrors, ok := e.(interface{ Unwrap() []error })
		if ok {
			for _, unwrapErr := range unwrappedErrors.Unwrap() {
				unwrap(unwrapErr)
			}
		} else {
			unwrapErr := errors.Unwrap(e)
			if unwrapErr != nil {
				unwrap(unwrapErr)
			}
		}
		errs = append(errs, e)
	}

	unwrap(err)

	return errs
}

// UnwrapErrorx returns the Errorx from an error.
func UnwrapErrorx(err error) Errorx {
	if err == nil {
		return nil
	}

	if !IsErrorx(err) {
		return nil
	}

	return NewFromError(err)
}

// UnwrapException returns the Exception from an Errorx.
func UnwrapException(err error) exceptionx.Exception {
	if !IsErrorx(err) {
		return nil
	}

	errx, ok := err.(Errorx)
	if ok && errx.Exception() != nil {
		return errx.Exception()
	}

	errx = UnwrapErrorx(err)
	if errx == nil {
		return nil
	}

	return errx.Exception()
}

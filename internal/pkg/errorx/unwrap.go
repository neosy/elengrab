package errorx

import (
	"errors"
	"slices"

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
// Duplicates are not excluded
func UnwrapAll(err error) []error {
	var (
		errs   []error
		unwrap func(error)
	)

	unwrap = func(e error) {
		if e == nil {
			return
		}

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

// UnwrapUnique error analysis errors in a slice
// Duplicates are excluded
func UnwrapUnique(err error) []error {
	errs := UnwrapAll(err)

	if len(errs) == 0 {
		return nil
	}

	if len(errs) == 1 {
		return []error{err}
	}

	errProcessed := make(map[string]struct{})
	newErrs := make([]error, 0, len(errs))

	for _, e := range errs {
		errStr := e.Error()
		if _, exists := errProcessed[errStr]; exists {
			continue
		}
		errProcessed[errStr] = struct{}{}

		newErrs = append(newErrs, e)
	}

	return newErrs
}

// findEdgeInChain is a helper function that traverses the error chain
// to find the first or last occurrence of an Errorx that satisfies a given condition.
func findEdgeInChain[T any](err error, reverse bool, fn func(Errorx) (T, bool)) (T, bool) {
	var zero T

	if err == nil {
		return zero, false
	}

	// If reverse is true, check the outermost error first before traversing the chain.
	if reverse {
		if errx, ok := err.(Errorx); ok {
			if v, ok := fn(errx); ok {
				return v, true
			}
		}
	}

	if !IsErrorx(err) {
		return zero, false
	}

	// Traverse the error chain and apply the provided function to each Errorx encountered.
	errs := UnwrapAll(err)
	if reverse {
		slices.Reverse(errs)
	}

	for _, err := range errs {
		if err == nil {
			continue
		}
		if errx, ok := err.(Errorx); ok {
			if v, ok := fn(errx); ok {
				return v, true
			}
		}
	}

	return zero, false
}

// OuterErrorx returns the first (outermost) Errorx found in the error chain.
func OuterErrorx(err error) Errorx {
	errx, _ := findEdgeInChain(err, true, func(e Errorx) (Errorx, bool) {
		return e, true
	})
	return errx
}

// RootErrorx returns the deepest (root cause) Errorx found in the error chain.
func RootErrorx(err error) Errorx {
	errx, _ := findEdgeInChain(err, false, func(e Errorx) (Errorx, bool) {
		return e, true
	})
	return errx
}

// OuterException returns the first (outermost) exception found in the error chain.
func OuterException(err error) exceptionx.Exception {
	ex, _ := findEdgeInChain(err, true, func(e Errorx) (exceptionx.Exception, bool) {
		ex := e.Exception()
		return ex, ex != nil
	})
	return ex
}

// RootException returns the deepest (root cause) exception found in the error chain.
func RootException(err error) exceptionx.Exception {
	ex, _ := findEdgeInChain(err, false, func(e Errorx) (exceptionx.Exception, bool) {
		ex := e.Exception()
		return ex, ex != nil
	})
	return ex
}

// OuterHttpStatus returns the first (outermost) HttpStatus found in the error chain.
func OuterHttpStatus(err error) int {
	status, _ := findEdgeInChain(err, true, func(e Errorx) (int, bool) {
		status := e.HttpStatus()
		return status, status != 0
	})
	return status
}

// RootHttpStatus returns the deepest (root cause) HttpStatus found in the error chain.
func RootHttpStatus(err error) int {
	status, _ := findEdgeInChain(err, false, func(e Errorx) (int, bool) {
		status := e.HttpStatus()
		return status, status != 0
	})
	return status
}

// OuterMessage returns the first (outermost) Message found in the error chain.
func OuterMessage(err error) string {
	msg, _ := findEdgeInChain(err, true, func(e Errorx) (string, bool) {
		msg := e.Message()
		return msg, msg != ""
	})
	return msg
}

// RootMessage returns the deepest (root cause) Message found in the error chain.
func RootMessage(err error) string {
	msg, _ := findEdgeInChain(err, false, func(e Errorx) (string, bool) {
		msg := e.Message()
		return msg, msg != ""
	})
	return msg
}

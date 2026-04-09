package errorx

import (
	"errors"
	"fmt"
	"sync"
)

// wrapErrorsDefaultSeparator is a constant string used as the default separator
const wrapErrorsDefaultSeparator = "; "

// wrapErrorsSeparator, separator for wrapErrors function
var (
	wrapErrorsSeparatorConfig = struct {
		separator string
		once      sync.Once
	}{
		separator: wrapErrorsDefaultSeparator,
	}
)

// SetWrapErrorsSeparator set separator for WrapErrors only once
// Subsequent calls to this function will have no effect.
// advice: call in main
func SetWrapErrorsSeparator(newSep string) {
	wrapErrorsSeparatorConfig.once.Do(func() {
		wrapErrorsSeparatorConfig.separator = newSep
	})
}

// wrapErrorsSeparator returns the separator used for wrapping error messages
func wrapErrorsSeparator() string {
	return wrapErrorsSeparatorConfig.separator
}

// WrapErrors combining errors into one
func WrapErrors(errs ...error) error {
	if len(errs) == 0 {
		return nil
	}

	var wrapped error

	for _, err := range errs {
		if err == nil {
			continue
		}

		// If wrapped is nil, set it to the current error and continue to the next iteration.
		if wrapped == nil {
			wrapped = err
			continue
		}

		// If an error of type errorx is encountered, convert wrapped to errorx if its type differs.
		// Then extract the error of type error from the structure in order to wrap it, instead of the error of type errorx.
		if e, ok := err.(errorxInternal); ok {
			if _, ok := wrapped.(errorxInternal); !ok {
				wrapped = NewFromError(wrapped)
			}
			err = e.Err()
		}

		switch e := wrapped.(type) {
		case errorxInternal:
			e.setErr(fmt.Errorf("%w%s%w", e.Err(), wrapErrorsSeparator(), err))
			wrapped = e.normalizeToInnerType()
		default:
			wrapped = fmt.Errorf("%w%s%w", wrapped, wrapErrorsSeparator(), err)
		}
	}

	return wrapped
}

// IsErrorx reports whether err is or wraps an Errorx.
func IsErrorx(err error) bool {
	if err == nil {
		return false
	}

	_, ok := err.(Errorx)
	if ok {
		return true
	}

	var e Errorx
	return errors.As(err, &e)
}

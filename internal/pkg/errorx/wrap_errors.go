package errorx

import (
	"errors"
	"fmt"
	"sync"
)

// wrapErrorsDefaultSeparator is a constant string used as the default separator
const wrapErrorsDefaultSeparator = ": "

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

		if wrapped == nil {
			wrapped = err
			continue
		}

		// Если встретилась ошибка errorx, то преобразуем wrapped в errorx, если его тип отличается
		// Далее из структуры извлекаем ошибку с типом error, чтобы ее обернуть, вместо ошибки с типом errorx
		if e, ok := err.(errorxInternal); ok {
			if _, ok := wrapped.(errorxInternal); !ok {
				wrapped = NewFromError(wrapped, e.args()...)
			} else {
				if cmb, ok := wrapped.(errorxInternal); ok {
					var args = make([]any, 0, 3)
					if cmb.Exception() == nil && e.Exception() != nil {
						args = append(args, e.Exception())
					}
					if cmb.Message() == nil && e.Message() != nil {
						text := e.Message()
						args = append(args, ErrorMessageArg(*text))
					}
					if cmb.HttpStatusCodeRaw() == nil && e.HttpStatusCodeRaw() != nil {
						code := e.HttpStatusCodeRaw()
						args = append(args, HttpStatusArg(*code))
					}
					cmb.initFromArgs(args...)
				}
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

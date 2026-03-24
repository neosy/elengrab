package errorx

import (
	"fmt"
	"sync"
)

// combineErrorsDefaultSeparator is a constant string used as the default separator
const combineErrorsDefaultSeparator = ": "

// combineErrorsSeparator, separator for CombineErrors function
var (
	combineErrorsSeparatorConfig = struct {
		separator string
		once      sync.Once
	}{
		separator: combineErrorsDefaultSeparator,
	}
)

// SetCombineErrorsSeparator set separator for CombineErrors only once
// Subsequent calls to this function will have no effect.
// advice: call in main
func SetCombineErrorsSeparator(newSep string) {
	combineErrorsSeparatorConfig.once.Do(func() {
		combineErrorsSeparatorConfig.separator = newSep
	})
}

// combineErrorsSeparator returns the separator used for combining error messages
func combineErrorsSeparator() string {
	return combineErrorsSeparatorConfig.separator
}

// CombineErrors combining errors into one
func CombineErrors(errs ...error) error {
	if len(errs) == 0 {
		return nil
	}

	var combined error

	for _, err := range errs {
		if err == nil {
			continue
		}

		if combined == nil {
			combined = err
			continue
		}

		// Если встретилась ошибка errorx, то преобразуем combined в errorx, если его тип отличается
		// Далее из структуры извлекаем ошибку с типом error, чтобы ее обернуть, вместо ошибки с типом errorx
		if e, ok := err.(Errorx); ok {
			if _, ok := combined.(Errorx); !ok {
				combined = NewFromError(combined, e.Args()...)
			} else {
				if cmb, ok := combined.(*errorx); ok {
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

		switch e := combined.(type) {
		case *errorx:
			e.err = fmt.Errorf("%w%s%w", e.err, combineErrorsSeparator(), err)
		default:
			combined = fmt.Errorf("%w%s%w", combined, combineErrorsSeparator(), err)
		}
	}

	return combined
}

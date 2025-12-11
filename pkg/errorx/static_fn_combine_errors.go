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

	var combined = errs[0]

	for i, err := range errs {
		// Пропускаем 1-ую, т.к. мы ее инициализировали в combined
		if i == 0 || err == nil {
			continue
		}

		// Создаем новый указатель если первая ошибка с типом Errorx
		if i == 1 {
			if _, ok := combined.(ErrorxInterface); ok {
				combined = NewByErr(combined)
			}
		}

		// Если встретилась ошибка Errorx, то преобразуем combined в Errorx, если его тип отличается
		// Далее из структуры извлекаем ошибку с типом error, чтобы ее обернуть, вместо ошибки с типом Errorx
		if e, ok := err.(ErrorxInterface); ok {
			if _, ok := combined.(ErrorxInterface); !ok {
				combined = NewByErr(
					combined,
					e.Message(),
					e.ExceptionType(),
					e.ExceptionCode(),
				)
			} else {
				// Если первый элемент был инициализирован из обычной ошибки
				if combined.(ErrorxInterface).ExceptionType() == nil && combined.(ErrorxInterface).ExceptionCode() == nil {
					if cmb, ok := combined.(*Errorx); ok {
						cmb.initFromArgs(e.ExceptionType(), e.ExceptionCode())
					}
				}
			}

			err = e.Err()
		}

		switch e := combined.(type) {
		case *Errorx:
			e.err = fmt.Errorf("%w%s%w", e.err, combineErrorsSeparator(), err)
		default:
			combined = fmt.Errorf("%w%s%w", combined, combineErrorsSeparator(), err)
		}
	}

	return combined
}

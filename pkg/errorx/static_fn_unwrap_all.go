package errorx

import (
	"errors"
)

// UnwrapAll error analysis errors in a slice
// Duplicates are excluded
func UnwrapAll(err error) (errs []error) {
	errProcessed := make(map[string]struct{})

	var unwrap func(error)
	unwrap = func(e error) {
		if e == nil {
			return
		}

		errStr := e.Error()
		if _, exists := errProcessed[errStr]; exists {
			return
		}

		errProcessed[errStr] = struct{}{}

		// Проверка на наличие нескольких вложенных ошибок (для errors.Join)
		unwrappedErrors, ok := e.(interface{ Unwrap() []error })
		if ok {
			for _, ue := range unwrappedErrors.Unwrap() {
				unwrap(ue)
			}
		} else {
			ue := errors.Unwrap(e)
			if ue != nil {
				unwrap(ue)
			} else {
				// Добавляем конечную ошибку в список
				errs = append(errs, e)
			}
		}
	}

	unwrap(err)

	return errs
}

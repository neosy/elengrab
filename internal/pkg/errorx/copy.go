package errorx

import (
	"errors"
	"strings"
)

// Copy copying an error including nested
func Copy(err error) error {
	if err == nil {
		return nil
	}

	isMultiLine := func(err error, firstErr error) bool {
		if err == nil || firstErr == nil {
			return false
		}
		msg := err.Error()
		msg = strings.TrimPrefix(msg, firstErr.Error())

		return strings.HasPrefix(msg, "\n")
	}

	var unwrap func(err error) error
	unwrap = func(err error) error {
		switch e := err.(type) {
		case (interface{ Unwrap() []error }):
			var errs = make([]error, 0, len(e.Unwrap()))
			for _, subErr := range e.Unwrap() {
				errs = append(errs, unwrap(subErr))
			}
			if len(errs) > 0 && !isMultiLine(err, errs[0]) {
				return WrapErrors(errs...)

			}
			return errors.Join(errs...)
		case errorxInternal:
			errx := New(e.Error(), e.args()...).(errorxInternal)
			errx.setErr(Copy(e.Err()))
			return errx.normalizeToInnerType()
		default:
			return e
		}
	}

	return unwrap(err)
}

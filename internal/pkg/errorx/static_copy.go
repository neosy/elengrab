package errorx

import (
	"errors"
)

// Copy copying an error including nested
func Copy(err error) error {
	if err == nil {
		return nil
	}

	var unwrap func(err error) error

	unwrap = func(err error) error {
		switch e := err.(type) {
		case (interface{ Unwrap() []error }):
			var errs = make([]error, 0, len(e.Unwrap()))
			for _, subErr := range e.Unwrap() {
				errs = append(errs, unwrap(subErr))
			}
			return errors.Join(errs...)
		case Errorx:
			errx := New(e.Error(), e.Args()...)
			errx.(*errorx).err = Copy(e.Err())
			return errx
		default:
			return e
		}
	}

	unwrap(err)

	return err
}

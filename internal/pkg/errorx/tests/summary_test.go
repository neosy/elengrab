package tests

import (
	"errors"
	"fmt"
	"testing"

	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

func TestErrorScenarios(t *testing.T) {
	err1 := fmt.Errorf("error 1")
	err2 := fmt.Errorf("error 2")
	errx2 := errorx.New("error 2", exceptionx.FORBIDDEN)

	// Таблица сценариев
	tests := []struct {
		name string
		fn   func() string
		want string
	}{
		{
			name: "fmt.Errorf %w: %w chain",
			fn: func() string {
				err := fmt.Errorf("%w: %w", err1, err2)
				err = fmt.Errorf("%w: %w", err, errors.New("error 3"))
				return err.Error()
			},
			want: "error 1: error 2: error 3",
		},
		{
			name: "errorx.WrapErrors",
			fn: func() string {
				err := fmt.Errorf("error 1")
				err = errorx.WrapErrors(err, fmt.Errorf("error 2"), fmt.Errorf("error 3"))
				return err.Error()
			},
			want: "error 1; error 2; error 3",
		},
		{
			name: "errors.Join",
			fn: func() string {
				err := errors.Join(err1, err2)
				return err.Error()
			},
			want: "error 1\nerror 2",
		},
		{
			name: "errorx.Copy",
			fn: func() string {
				err := errors.Join(err1, err2)
				copied := errorx.Copy(err)
				return copied.Error()
			},
			want: "error 1\nerror 2",
		},
		{
			name: "errorx.Wrap multiple",
			fn: func() string {
				err := fmt.Errorf("error 1")
				err = fmt.Errorf("error 2: %w", err)
				err = errorx.New("error 3", exceptionx.ERROR, err)
				err = err.(errorx.Errorx).Wrap(errors.New("error 3.1"), errors.New("error 3.2"))
				return err.Error()
			},
			want: "error 3; error 2: error 1; error 3.1; error 3.2",
		},
		{
			name: "errorx.Wrap and errorx.Join",
			fn: func() string {
				err := fmt.Errorf("error 1")
				err = fmt.Errorf("error 2: %w", err)
				err = errorx.New("error 3", exceptionx.ERROR, err)
				err = err.(errorx.Errorx).Wrap(errors.New("error 3.1"), errors.New("error 3.2"))
				err = errors.Join(errors.New("error 4"), err)
				return fmt.Sprintf("%v %v", err.Error(), errorx.OuterException(err))
			},
			want: "error 4\nerror 3; error 2: error 1; error 3.1; error 3.2 internal Server Error",
		},
		{
			name: "fmt.Errorf, errorx.Join, DomainException and OuterException",
			fn: func() string {
				err := fmt.Errorf("error 1")
				err = fmt.Errorf("error 2: %w", err)
				err = errorx.New("error 3", exceptionx.ERROR, err)
				err = err.(errorx.Errorx).Wrap(errors.New("error 3.1"), errors.New("error 3.2"))
				err = errors.Join(errors.New("error 4"), err)
				err = errorx.NewFromError(err, exceptionx.NOT_FOUND).Join(fmt.Errorf("error 5"))
				return fmt.Sprintf("%v %v", err.Error(), errorx.OuterException(err))
			},
			want: "error 4\nerror 3; error 2: error 1; error 3.1; error 3.2\nerror 5 not Found",
		},
		{
			name: "fmt.Errorf, errorx.Join, errorx.Wrap, DomainException and OuterException",
			fn: func() string {
				err := fmt.Errorf("error 1")
				err = fmt.Errorf("error 2: %w", err)
				err = errorx.New("error 3", exceptionx.ERROR, err)
				err = err.(errorx.Errorx).Wrap(errors.New("error 3.1"), errors.New("error 3.2"))
				err = errors.Join(errors.New("error 4"), err)
				err = errorx.NewFromError(err, exceptionx.NOT_FOUND).Join(fmt.Errorf("error 5"))
				err = err.(errorx.Errorx).Wrap(errors.New("error 6"))
				err = err.(errorx.Errorx).Wrap(errors.New("error 7"))
				return fmt.Sprint(err.Error())
			},
			want: "error 4\nerror 3; error 2: error 1; error 3.1; error 3.2\nerror 5; error 6; error 7",
		},
		{
			name: "Copy [fmt.Errorf, errorx.Join, errorx.Wrap, DomainException and OuterException]",
			fn: func() string {
				err := fmt.Errorf("error 1")
				err = fmt.Errorf("error 2: %w", err)
				err = errorx.New("error 3", exceptionx.ERROR, err)
				err = err.(errorx.Errorx).Wrap(errors.New("error 3.1"), errors.New("error 3.2"))
				err = errors.Join(errors.New("error 4"), err)
				err = errorx.NewFromError(err, exceptionx.NOT_FOUND).Join(fmt.Errorf("error 5"))
				err = err.(errorx.Errorx).Wrap(errors.New("error 6"))
				err = err.(errorx.Errorx).Wrap(errors.New("error 7"))
				return fmt.Sprint(err.(errorx.Errorx).Copy())
			},
			want: "error 4\nerror 3; error 2: error 1; error 3.1; error 3.2\nerror 5; error 6; error 7",
		},
		{
			name: "fmt.Errorf(wrapped), exception Code and Message",
			fn: func() string {
				err := fmt.Errorf("%w: %w", err1, errx2)
				return fmt.Sprintf("%v %v %v", err.Error(), errorx.OuterException(err).String(), errorx.OuterException(err).Code())
			},
			want: "error 1: error 2 Forbidden FORBIDDEN",
		},
		{
			name: "UnwrapAll()",
			fn: func() string {
				err := fmt.Errorf("%w: %w", err1, errx2)
				err = fmt.Errorf("error 3: %w", err)
				err = fmt.Errorf("error 4: %w", err)
				err = errorx.NewFromError(err, exceptionx.ERROR)
				return fmt.Sprint(errorx.UnwrapAll(err))
			},
			want: "[error 1 error 2 error 1: error 2 error 3: error 1: error 2 error 4: error 3: error 1: error 2]",
		},
		{
			name: "err.UnwrapTexts().Join()",
			fn: func() string {
				err := fmt.Errorf("%w: %w", err1, errx2)
				err = fmt.Errorf("error 3: %w", err)
				err = fmt.Errorf("error 4: %w", err)
				err = errorx.NewFromError(err, exceptionx.ERROR)
				return fmt.Sprint(err.(errorx.Errorx).UnwrapTexts().Join())
			},
			want: "error 1; error 2; error 1: error 2; error 3: error 1: error 2; error 4: error 3: error 1: error 2",
		},
		{
			name: "errorx.NewErrorTexts().AddUnwrapErr(err).Join()",
			fn: func() string {
				err := fmt.Errorf("%w: %w", err1, errx2)
				err = fmt.Errorf("error 3: %w", err)
				err = fmt.Errorf("error 4: %w", err)
				err = errorx.NewFromError(err, exceptionx.ERROR)
				return fmt.Sprint(errorx.NewErrorTexts().AddUnwrapErr(err).Join())
			},
			want: "error 1; error 2; error 1: error 2; error 3: error 1: error 2; error 4: error 3: error 1: error 2",
		},
		{
			name: "err.Error() and OuterException",
			fn: func() string {
				err := fmt.Errorf("%w: %w", err1, errx2)
				err = fmt.Errorf("error 3: %w", err)
				err = fmt.Errorf("error 4: %w", err)
				err = errorx.NewFromError(err, exceptionx.ERROR)
				return fmt.Sprintf("%v %v", err.Error(), errorx.OuterException(err))
			},
			want: "error 4: error 3: error 1: error 2 internal Server Error",
		},
		{
			name: "errors.Unwrap(err)",
			fn: func() string {
				err := fmt.Errorf("%w: %w", err1, errx2)
				err = fmt.Errorf("error 3: %w", err)
				err = fmt.Errorf("error 4: %w", err)
				err = errorx.NewFromError(err, exceptionx.ERROR)
				return fmt.Sprint(errors.Unwrap(err).Error())
			},
			want: "error 3: error 1: error 2",
		},
		{
			name: "errors.Unwrap(errors.Unwrap(err))",
			fn: func() string {
				err := fmt.Errorf("%w: %w", err1, errx2)
				err = fmt.Errorf("error 3: %w", err)
				err = fmt.Errorf("error 4: %w", err)
				err = errorx.NewFromError(err, exceptionx.ERROR)
				return fmt.Sprint(errors.Unwrap(errors.Unwrap(err)).Error())
			},
			want: "error 1: error 2",
		},
		{
			name: "errors.Is(err, err1)",
			fn: func() string {
				err := fmt.Errorf("%w: %w", err1, errx2)
				err = fmt.Errorf("error 3: %w", err)
				err = fmt.Errorf("error 4: %w", err)
				err = errorx.NewFromError(err, exceptionx.ERROR)
				if errors.Is(err, err1) {
					return fmt.Sprint(err1.Error())
				}
				return ""
			},
			want: "error 1",
		},
		{
			name: "errors.Is(err, errx2)",
			fn: func() string {
				err := fmt.Errorf("%w: %w", err1, errx2)
				err = fmt.Errorf("error 3: %w", err)
				err = fmt.Errorf("error 4: %w", err)
				err = errorx.NewFromError(err, exceptionx.ERROR)
				if errors.Is(err, errx2) {
					return fmt.Sprintf("%v %v", err2, errorx.OuterException(err))
				}
				return ""
			},
			want: "error 2 internal Server Error",
		},
		{
			name: "errorx.Copy(err), Exception code, Exception message",
			fn: func() string {
				err := fmt.Errorf("%w: %w", err1, errx2)
				err = fmt.Errorf("error 3: %w", err)
				err = fmt.Errorf("error 4: %w", err)
				err = errorx.NewFromError(err, exceptionx.ERROR)
				newErr := err.(errorx.Errorx).Copy()
				return fmt.Sprintf("%v %v %v", newErr.Error(), errorx.OuterException(newErr).Code(), errorx.OuterException(newErr).Message())
			},
			want: "error 4: error 3: error 1: error 2 ERROR Internal Server Error",
		},
		{
			name: "errorx.Errorf(), OuterException, HttpStatusCode",
			fn: func() string {
				err := errorx.Errorf("error 2: %w", err1, exceptionx.NOT_FOUND)
				err = errorx.Errorf("error %d: %w", 3, err)
				return fmt.Sprintf("%v %v %v %v", err.Error(), errorx.OuterException(err).Code(), errorx.OuterException(err), err.HttpStatusCode())
			},
			want: "error 3: error 2: error 1 NOT_FOUND not Found 404",
		},
	}

	// Запуск тестов
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fn()
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

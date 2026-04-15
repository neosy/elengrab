package ierrors

import (
	"github.com/neosy/elengrab/internal/exceptions"
	"github.com/neosy/elengrab/internal/pkg/errorx"
)

var (
	ErrFuncParamNullPointer = exceptions.FUNCTION_PARAMETER_NULL_POINTER.NewErrorx(
		errorx.WithErrorMessage("Something went wrong"),
	).(errorx.Errorx)
	ErrFileNotFound = exceptions.FILE_NOT_FOUND.NewErrorx(
		errorx.WithErrorMessage("The requested file could not be found"),
	).(errorx.Errorx)
)

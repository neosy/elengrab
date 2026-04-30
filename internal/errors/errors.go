package ierrors

import (
	"github.com/neosy/elengrab/internal/exceptions"
	"github.com/neosy/elengrab/internal/pkg/errorx"
)

var (
	ErrFuncParamNullPointer = exceptions.FUNCTION_PARAMETER_NULL_POINTER.NewErrorx(
		errorx.WithErrorMessage("Something went wrong"),
	).(errorx.Errorx)
	ErrFuncContainsEmptyFields = exceptions.FUNCTION_CONTAINS_EMPTY_FIELDS.NewErrorx(
		errorx.WithErrorMessage("Something went wrong"),
	).(errorx.Errorx)
	ErrFileNotFound = exceptions.FILE_NOT_FOUND.NewErrorx(
		errorx.WithErrorMessage("The requested file could not be found"),
	).(errorx.Errorx)
	ErrThumbnailNotFound = exceptions.THUMBNAIL_NOT_FOUND.NewErrorx(
		errorx.WithErrorMessage("The requested thumbnail could not be found"),
	).(errorx.Errorx)
)

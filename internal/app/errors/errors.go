package apperrors

import (
	ierrors "github.com/neosy/elengrab/internal/errors"
	"github.com/neosy/elengrab/internal/exceptions"
	"github.com/neosy/elengrab/internal/pkg/errorx"
)

var (
	ErrFuncParamNullPointer    = ierrors.ErrFuncParamNullPointer
	ErrFuncContainsEmptyFields = ierrors.ErrFuncContainsEmptyFields
	ErrDownloaderEmptyResponse = exceptions.DOWNLOADER_EMPTY_RESPONSE.NewErrorx(
		errorx.WithErrorMessage("We couldn't retrieve the data. Please try again later"),
	).(errorx.Errorx)
	ErrFileNotFound = ierrors.ErrFileNotFound
	ErrFileIDIsNil  = exceptions.FILE_ID_IS_NIL.NewErrorx(
		errorx.WithErrorMessage("Invalid request"),
	).(errorx.Errorx)
)

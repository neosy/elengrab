package apperrors

import (
	ierrors "github.com/neosy/elengrab/internal/errors"
	"github.com/neosy/elengrab/internal/exceptions"
	"github.com/neosy/elengrab/internal/pkg/errorx"
)

var (
	ErrFuncParamNullPointer    = ierrors.ErrFuncParamNullPointer
	ErrDownloaderEmptyResponse = exceptions.DOWNLOADER_EMPTY_RESPONSE.NewErrorx(
		errorx.WithErrorMessage("We couldn't retrieve the data. Please try again later"),
	)
	ErrFileNotFound = exceptions.FILE_NOT_FOUND.NewErrorx(
		errorx.WithErrorMessage("The requested file could not be found"),
	)
	ErrFileIDIsNil = exceptions.FILE_ID_IS_NIL.NewErrorx(
		errorx.WithErrorMessage("Invalid request"),
	)
)

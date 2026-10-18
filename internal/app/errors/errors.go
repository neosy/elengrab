package apperrors

import (
	"errors"

	ierrors "github.com/neosy/elengrab/internal/errors"
	"github.com/neosy/elengrab/internal/exceptions"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/valyala/fasthttp"
)

var (
	ErrFuncParamNullPointer    = ierrors.ErrFuncParamNullPointer
	ErrFuncContainsEmptyFields = ierrors.ErrFuncContainsEmptyFields
	ErrDownloaderEmptyResponse = exceptions.DOWNLOADER_EMPTY_RESPONSE.NewErrorx(
		errorx.WithErrorMessage("We couldn't retrieve the data. Please try again later"),
	).(errorx.Errorx)
	ErrDownloadNotFound = ierrors.ErrDownloadNotFound
	ErrDownloadIDIsNil  = exceptions.DOWNLOAD_ID_IS_NIL.NewErrorx(
		errorx.WithErrorMessage("Invalid request"),
	).(errorx.Errorx)

	ErrIDConflict = errors.New("id conflict")

	ErrUnsupportedMediaFormat = errorx.New(
		"unsupported media format",
		errorx.WithErrorMessage("Unsupported media format"),
		errorx.WithHttpStatus(fasthttp.StatusUnsupportedMediaType),
	)
	ErrUnsupportedVideoCodec = errorx.New(
		"unsupported video codec",
		errorx.WithErrorMessage("Unsupported video codec"),
		errorx.WithHttpStatus(fasthttp.StatusUnsupportedMediaType),
	)
	ErrUnsupportedAudioCodec = errorx.New(
		"unsupported audio codec",
		errorx.WithErrorMessage("Unsupported audio codec"),
		errorx.WithHttpStatus(fasthttp.StatusUnsupportedMediaType),
	)
)

package apierrors

import (
	ierrors "github.com/neosy/elengrab/internal/errors"
	"github.com/neosy/elengrab/internal/exceptions"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/valyala/fasthttp"
)

var (
	ErrQueryReturnedEmptyResult = exceptions.QUERY_RETURNED_EMPTY_RESULT.NewErrorx(
		errorx.WithErrorMessage("Something went wrong"),
	).(errorx.Errorx)

	ErrEmptyResponse = exceptions.EMPTY_RESPONSE.NewErrorx(
		errorx.WithErrorMessage("Something went wrong"),
	).(errorx.Errorx)

	ErrHTTPSRequired = errorx.NewHTTPMessage("HTTPS is required", fasthttp.StatusUpgradeRequired)

	ErrFileNotFound          = ierrors.ErrFileNotFound
	ErrDownloadIDIsRequired  = errorx.NewMessage("downloadId is required", exceptions.INVALID_REQUEST)
	ErrDownloadIDIsIncorrect = errorx.NewMessage("downloadId is incorrect", exceptions.INVALID_REQUEST)

	ErrThumbnailNotFound      = ierrors.ErrThumbnailNotFound
	ErrThumbnailIdIsRequired  = errorx.NewMessage("thumbnailId is required", exceptions.INVALID_REQUEST)
	ErrThumbnailIdIsIncorrect = errorx.NewMessage("thumbnailId is incorrect", exceptions.INVALID_REQUEST)

	ErrChannelIsRequired  = errorx.NewMessage("channelId is required", exceptions.INVALID_REQUEST)
	ErrChannelIsIncorrect = errorx.NewMessage("channelId is incorrect", exceptions.INVALID_REQUEST)

	ErrURLIsRequired = errorx.New("url is required", exceptions.INVALID_REQUEST, errorx.WithErrorMessage("URL is required"))
	ErrInvalidURL    = errorx.NewMessage("invalid URL", exceptions.INVALID_REQUEST)

	ErrUserIDIsRequired  = errorx.NewMessage("userId is required", exceptions.INVALID_REQUEST)
	ErrUserIDIsIncorrect = errorx.NewMessage("uerId is incorrect", exceptions.INVALID_REQUEST)
)

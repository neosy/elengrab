package apierrors

import (
	ierrors "github.com/neosy/elengrab/internal/errors"
	"github.com/neosy/elengrab/internal/exceptions"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

var (
	ErrQueryReturnedEmptyResult = exceptions.QUERY_RETURNED_EMPTY_RESULT.NewErrorx(
		errorx.WithErrorMessage("Something went wrong"),
	).(errorx.Errorx)

	ErrEmptyResponse = exceptions.EMPTY_RESPONSE.NewErrorx(
		errorx.WithErrorMessage("Something went wrong"),
	).(errorx.Errorx)

	ErrUnauthorized = exceptionx.UNAUTHORIZED.NewErrorx().(errorx.Errorx)

	ErrFileNotFound      = ierrors.ErrFileNotFound
	ErrFileIdIsRequired  = errorx.NewWithMessage("fileId is required", exceptions.INVALID_REQUEST)
	ErrFileIdIsIncorrect = errorx.NewWithMessage("fileId is incorrect", exceptions.INVALID_REQUEST)

	ErrThumbnailNotFound      = ierrors.ErrThumbnailNotFound
	ErrThumbnailIdIsRequired  = errorx.NewWithMessage("thumbnailId is required", exceptions.INVALID_REQUEST)
	ErrThumbnailIdIsIncorrect = errorx.NewWithMessage("thumbnailId is incorrect", exceptions.INVALID_REQUEST)

	ErrChannelIsRequired  = errorx.NewWithMessage("channelId is required", exceptions.INVALID_REQUEST)
	ErrChannelIsIncorrect = errorx.NewWithMessage("channelId is incorrect", exceptions.INVALID_REQUEST)

	ErrURLIsRequired = errorx.NewWithMessage("URL is required", exceptions.INVALID_REQUEST)
	ErrInvalidURL    = errorx.NewWithMessage("invalid URL", exceptions.INVALID_REQUEST)
)

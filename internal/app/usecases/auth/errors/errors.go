package autherr

import (
	"errors"

	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

var (
	ErrSessionExpired  = errors.New("session expired")
	ErrSessionNotFound = errorx.New("session not found", exceptionx.NOT_FOUND)

	ErrUserNotFound    = errorx.New("user not found", exceptionx.NOT_FOUND)
	ErrUserIsNotActive = errorx.New("user is not active", exceptionx.UNAUTHORIZED)
	ErrUserDeleted     = errorx.New("user deleted", exceptionx.UNAUTHORIZED)

	ErrInternal             = errorx.New("internal server error", exceptionx.ERROR)
	ErrFunctionNilParameter = errorx.New("function parameter is nil", exceptionx.ERROR)
)

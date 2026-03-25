package auth

import (
	"errors"

	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

var (
	ErrSessionExpired  = errors.New("session expired")
	ErrSessionNotFound = errorx.New("session not found", exceptionx.NOT_FOUND)
	ErrUserNotFound    = errorx.New("user not found", exceptionx.NOT_FOUND)
	ErrInternal        = errorx.New("internal server error", exceptionx.ERROR)
)

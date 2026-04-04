package linkerr

import (
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

var (
	ErrLinkNotFound = errorx.New("link not found", exceptionx.NOT_FOUND)

	ErrInternal             = errorx.New("internal server error", exceptionx.ERROR)
	ErrFunctionNilParameter = errorx.New("function parameter is nil", exceptionx.ERROR)
)

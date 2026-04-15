package ierrors

import (
	"github.com/neosy/elengrab/internal/exceptions"
	"github.com/neosy/elengrab/internal/pkg/errorx"
)

var (
	ErrFuncParamNullPointer = exceptions.FUNCTION_PARAMETER_NULL_POINTER.NewErrorx(
		errorx.WithErrorMessage("Something went wrong"),
	)
)

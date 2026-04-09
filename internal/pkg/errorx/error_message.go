package errorx

import "github.com/neosy/elengrab/internal/pkg/errorx/internal/types"

// ErrorMessageProvider returns an error message string.
type ErrorMessageProvider = types.ErrorMessageProvider

// ErrorMessageArg creates a provider that always returns the given text.
var ErrorMessageArg = types.ErrorMessageArg

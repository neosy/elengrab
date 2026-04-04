package errorx

import errmsg "github.com/neosy/elengrab/internal/pkg/errorx/internal/error_message"

// ErrorMessageProvider returns an error message string.
type ErrorMessageProvider = errmsg.ErrorMessageProvider

// ErrorMessageArg creates a provider that always returns the given text.
var ErrorMessageArg = errmsg.ErrorMessageArg

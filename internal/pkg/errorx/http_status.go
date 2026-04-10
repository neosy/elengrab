package errorx

import "github.com/neosy/elengrab/internal/pkg/errorx/internal/types"

type (
	// HttpStatusProvider is a function that provides an HTTP status.
	HttpStatusProvider = types.HttpStatusProvider
)

var (
	// HttpStatusArg returns a function that provides the HTTP status.
	WithHttpStatus = types.WithHttpStatus
)

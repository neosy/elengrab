package authz

import (
	"log/slog"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type Authorization struct {
	logger *slog.Logger

	// options
	appMode dtypes.AppMode
}

func NewAuthorization(
	logger *slog.Logger,

	// options
	appMode dtypes.AppMode,
) *Authorization {
	return &Authorization{
		logger: logger,

		// options
		appMode: appMode,
	}
}

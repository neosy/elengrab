package authz

import (
	"log/slog"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type Authorization struct {
	logger *slog.Logger

	// options
	historyMode dtypes.HistoryMode
}

func NewAuthorization(
	logger *slog.Logger,

	// options
	historyMode dtypes.HistoryMode,
) *Authorization {
	return &Authorization{
		logger: logger,

		// options
		historyMode: historyMode,
	}
}

package admin

import (
	"log/slog"

	pservices "github.com/neosy/elengrab/internal/ports/services"
)

type admin struct {
	logger *slog.Logger

	// services
	adminPanel pservices.AdminPanelService

	// options
}

func NewAdmin(
	logger *slog.Logger,

	// services
	adminPanel pservices.AdminPanelService,

	// options
) Admin {
	return &admin{
		logger: logger,

		// adminPanel
		adminPanel: adminPanel,

		// options
	}
}

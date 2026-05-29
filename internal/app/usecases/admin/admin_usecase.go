package admin

import (
	"log/slog"

	pservices "github.com/neosy/elengrab/internal/ports/services"
)

type Admin struct {
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
) *Admin {
	return &Admin{
		logger: logger,

		// adminPanel
		adminPanel: adminPanel,

		// options
	}
}

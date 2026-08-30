package link

import (
	"log/slog"

	ilink "github.com/neosy/elengrab/internal/app/usecases/link/internal/link"
	linkclick "github.com/neosy/elengrab/internal/app/usecases/link/internal/link_click"
	"github.com/neosy/elengrab/internal/app/usecases/link/mappers"
	"github.com/neosy/elengrab/internal/ports/persistence"
)

type link struct {
	logger  *slog.Logger
	mappers *mappers.Mappers

	// internal
	link      *ilink.Link
	linkClick *linkclick.LinkClick

	// options
	options LinkOptions
}

func NewLink(
	logger *slog.Logger,

	// repositories
	linkRepo persistence.LinkRepositoryFactory,
	linkClickRepo persistence.LinkClickRepositoryFactory,

	// options
	opts ...LinkOption,
) *link {
	options := DefaultLinkOptions()

	for _, opt := range opts {
		opt(&options)
	}

	link := &link{
		logger:  logger,
		mappers: mappers.NewMappers(),

		// internal
		link:      ilink.NewLink(logger, linkRepo),
		linkClick: linkclick.NewLinkClick(logger, linkClickRepo),

		// options
		options: options,
	}

	return link
}

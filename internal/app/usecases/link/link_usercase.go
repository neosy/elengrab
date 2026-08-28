package link

import (
	"log/slog"

	link "github.com/neosy/elengrab/internal/app/usecases/link/internal/link"
	linkclick "github.com/neosy/elengrab/internal/app/usecases/link/internal/link_click"
	"github.com/neosy/elengrab/internal/app/usecases/link/mappers"
	"github.com/neosy/elengrab/internal/ports/persistence"
)

type Link struct {
	logger  *slog.Logger
	mappers *mappers.Mappers

	// internal
	link      *link.Link
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
) *Link {
	options := DefaultLinkOptions()

	for _, opt := range opts {
		opt(&options)
	}

	link := &Link{
		logger:  logger,
		mappers: mappers.NewMappers(),

		// internal
		link:      link.NewLink(logger, linkRepo),
		linkClick: linkclick.NewLinkClick(logger, linkClickRepo),

		// options
		options: options,
	}

	return link
}

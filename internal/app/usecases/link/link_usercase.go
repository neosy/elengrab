package link

import (
	"log/slog"

	linklink "github.com/neosy/elengrab/internal/app/usecases/link/internal/link"
	linkclick "github.com/neosy/elengrab/internal/app/usecases/link/internal/link_click"
	"github.com/neosy/elengrab/internal/app/usecases/link/mappers"
	"github.com/neosy/elengrab/internal/ports/persistence"
)

type Link struct {
	logger  *slog.Logger
	mappers *mappers.Mappers

	// internal
	link      *linklink.Link
	linkClick *linkclick.LinkClick

	// options
	options LinkOptions
}

func NewLink(
	logger *slog.Logger,

	// repositories
	linkRep persistence.LinkRepository,
	linkClickRep persistence.LinkClickRepository,

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
		link:      linklink.NewLink(logger, linkRep),
		linkClick: linkclick.NewLinkClick(logger, linkClickRep),

		// options
		options: options,
	}

	return link
}

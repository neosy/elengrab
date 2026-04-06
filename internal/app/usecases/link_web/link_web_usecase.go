package linkweb

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/app/usecases/link"
	pservices "github.com/neosy/elengrab/internal/ports/services"
)

type LinkWeb struct {
	logger *slog.Logger

	// services
	link pservices.ShortLinkService

	// options
	shortCodeLength uint8
}

func NewLinkWeb(
	logger *slog.Logger,

	// services
	link *link.Link,

	// options
	shortCodeLength uint8,
) *LinkWeb {
	return &LinkWeb{
		logger: logger,

		// services
		link: link,

		// options
		shortCodeLength: shortCodeLength,
	}
}

package linkweb

import (
	"log/slog"
	"time"

	"github.com/neosy/elengrab/internal/app/usecases/link"
	pservices "github.com/neosy/elengrab/internal/ports/services"
)

type Dependencies struct {
	// Services
	Link *link.Link
}

type Options struct {
	BaseShortURL     string
	ShortCodeLength  uint8
	LinkTTL          time.Duration
	RefreshThreshold time.Duration
}

type LinkWeb struct {
	logger *slog.Logger

	// services
	link pservices.ShortLinkService

	// options
	options Options
}

func NewLinkWeb(logger *slog.Logger, deps Dependencies, opts Options) *LinkWeb {
	return &LinkWeb{
		logger: logger,

		// services
		link: deps.Link,

		// options
		options: opts,
	}
}

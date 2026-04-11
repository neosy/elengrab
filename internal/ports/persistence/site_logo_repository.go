package persistence

import (
	"context"

	"github.com/google/uuid"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
	memsimple "github.com/neosy/elengrab/internal/pkg/cache/memory/simple"
)

type SiteLogoRepository interface {
	Insert(ctx context.Context, logo *dmedia.SiteLogo) error
	Update(ctx context.Context, logo *dmedia.SiteLogo) error
	Save(ctx context.Context, logo *dmedia.SiteLogo) error
	FindByLogoID(ctx context.Context, logoID uuid.UUID) (*dmedia.SiteLogo, error)
	ExistsByLogoID(ctx context.Context, logoID uuid.UUID) (bool, error)
	FindBySiteURL(ctx context.Context, siteURL string) (*dmedia.SiteLogo, error)
	ExistsBySiteURL(ctx context.Context, siteURL string) (bool, error)
}

type SiteLogoCacheRepository interface {
	Save(ctx context.Context, logo *dmedia.SiteLogo) error
	SaveNegative(ctx context.Context, siteURL string) error
	FindByLogoID(ctx context.Context, logoID uuid.UUID) (*dmedia.SiteLogo, error)
	ExistsByLogoID(ctx context.Context, logoID uuid.UUID) (bool, error)
	FindBySiteURL(ctx context.Context, siteURL string) (*dmedia.SiteLogo, memsimple.CacheStatus, error)
	ExistsBySiteURL(ctx context.Context, siteURL string) (bool, error)
	CleanExpired(ctx context.Context) error
}

package inmemory

import (
	"context"
	"time"

	"github.com/google/uuid"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
	memsimple "github.com/neosy/elengrab/internal/pkg/cache/memory/simple"
)

type SiteLogoRepository struct {
	memsimple.Repository[dmedia.SiteLogo]

	cacheByLogoID  memsimple.Cache[uuid.UUID, dmedia.SiteLogo]
	cacheBySiteURL memsimple.Cache[string, dmedia.SiteLogo]
}

// newSiteLogoRepository returns a new object for the repository
func newSiteLogoRepository(ttl time.Duration) *SiteLogoRepository {
	r := &SiteLogoRepository{
		cacheByLogoID:  memsimple.NewCacheWithDeaultCopier[uuid.UUID, dmedia.SiteLogo, *dmedia.SiteLogo](),
		cacheBySiteURL: memsimple.NewCacheWithDeaultCopier[string, dmedia.SiteLogo, *dmedia.SiteLogo](),
	}
	r.Repository.Init(ttl)
	return r
}

func (r *SiteLogoRepository) Name() string {
	return "site_logo"
}

// Save saves a new logo to the repository.
func (r *SiteLogoRepository) Save(ctx context.Context, logo *dmedia.SiteLogo) error {
	if logo == nil || logo.LogoID == uuid.Nil || logo.SiteURL == "" {
		return nil
	}

	save := func() error {
		r.cacheByLogoID.Save(
			logo.LogoID,
			logo,
			r.TTL(),
		)
		r.cacheBySiteURL.Save(
			logo.SiteURL,
			logo,
			r.TTL(),
		)
		return nil
	}

	return r.Repository.Save(ctx, save)
}

// SaveNegative saves a negative entry for a given site URL to the repository.
func (r *SiteLogoRepository) SaveNegative(ctx context.Context, siteURL string) error {
	if siteURL == "" {
		return nil
	}

	save := func() error {
		r.cacheBySiteURL.Save(
			siteURL,
			nil,
			r.TTL(),
		)
		return nil
	}

	return r.Repository.Save(ctx, save)
}

func (r *SiteLogoRepository) Delete(ctx context.Context, logoID uuid.UUID, siteURL string) error {
	fnDelete := func() error {
		if logoID != uuid.Nil {
			r.cacheByLogoID.Delete(logoID)
		}
		if siteURL != "" {
			r.cacheBySiteURL.Delete(siteURL)
		}
		return nil
	}
	return r.Repository.Delete(ctx, fnDelete)
}

// FindByLogoID retrieves a site logo by its unique ID from the repository.
func (r *SiteLogoRepository) FindByLogoID(ctx context.Context, logoID uuid.UUID) (*dmedia.SiteLogo, error) {
	fnFind := func() (*dmedia.SiteLogo, error) {
		return r.cacheByLogoID.Find(logoID), nil
	}
	return r.Repository.Find(ctx, fnFind)
}

// Checks if a site logo exists in the repository by its logo ID.
func (r *SiteLogoRepository) ExistsByLogoID(ctx context.Context, logoID uuid.UUID) (bool, error) {
	fnExists := func() (bool, error) {
		return r.cacheByLogoID.Exists(logoID), nil
	}
	return r.Repository.Exists(ctx, fnExists)
}

// FindBySiteURL retrieves a site logo by its URL from the repository.
func (r *SiteLogoRepository) FindBySiteURL(ctx context.Context, siteURL string) (*dmedia.SiteLogo, memsimple.CacheStatus, error) {
	// Define a function to find the site logo in the cache.
	fnFind := func() (*dmedia.SiteLogo, memsimple.CacheStatus, error) {
		data, status := r.cacheBySiteURL.FindWithStatus(siteURL)
		return data, status, nil
	}
	return r.Repository.FindWithStatus(ctx, fnFind)
}

// Checks if a site logo exists in the repository based on the provided site URL.
func (r *SiteLogoRepository) ExistsBySiteURL(ctx context.Context, siteURL string) (bool, error) {
	fnExists := func() (bool, error) {
		return r.cacheBySiteURL.Exists(siteURL), nil
	}
	return r.Repository.Exists(ctx, fnExists)
}

// CleanExpired cleans expired items from the cache and calls the provided clean function.
func (r *SiteLogoRepository) CleanExpired(ctx context.Context) error {
	clean := func() error {
		r.cacheByLogoID.CleanExpired()
		r.cacheBySiteURL.CleanExpired()
		return nil
	}
	return r.Repository.CleanExpired(ctx, clean)
}

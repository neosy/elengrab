package inmemory

import (
	"context"
	"time"

	"github.com/google/uuid"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
	nmemory "github.com/neosy/elengrab/pkg/ncache/memory"
)

type SiteLogoRepository struct {
	nmemory.Repository[dmedia.SiteLogo]

	cacheByLogoID  nmemory.Cache[uuid.UUID, dmedia.SiteLogo]
	cacheBySiteURL nmemory.Cache[string, dmedia.SiteLogo]
}

// newSiteLogoRepository returns a new object for the repository
func newSiteLogoRepository(ttl time.Duration) *SiteLogoRepository {
	r := &SiteLogoRepository{
		cacheByLogoID:  nmemory.NewCacheWithDeaultCopier[uuid.UUID, dmedia.SiteLogo, *dmedia.SiteLogo](),
		cacheBySiteURL: nmemory.NewCacheWithDeaultCopier[string, dmedia.SiteLogo, *dmedia.SiteLogo](),
	}
	r.Repository.Init(ttl)
	return r
}

func (r *SiteLogoRepository) Save(_ context.Context, logo *dmedia.SiteLogo) error {
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

	return r.Repository.Save(save)
}

func (r *SiteLogoRepository) Insert(ctx context.Context, logo *dmedia.SiteLogo) error {
	return r.Save(ctx, logo)
}

func (r *SiteLogoRepository) Update(ctx context.Context, logo *dmedia.SiteLogo) error {
	return r.Save(ctx, logo)
}

func (r *SiteLogoRepository) Delete(_ context.Context, logoID uuid.UUID) error {
	fnDelete := func() error {
		logo := r.cacheByLogoID.Find(logoID)
		if logo == nil {
			return nil
		}
		r.cacheByLogoID.Delete(logoID)
		r.cacheBySiteURL.Delete(logo.SiteURL)
		return nil
	}
	return r.Repository.Delete(fnDelete)
}

// FindByLogoID retrieves a site logo by its unique ID from the repository.
func (r *SiteLogoRepository) FindByLogoID(_ context.Context, logoID uuid.UUID) (*dmedia.SiteLogo, error) {
	fnFind := func() (*dmedia.SiteLogo, error) {
		return r.cacheByLogoID.Find(logoID), nil
	}
	return r.Repository.Find(fnFind)
}

// Checks if a site logo exists in the repository by its logo ID.
func (r *SiteLogoRepository) ExistsByLogoID(_ context.Context, logoID uuid.UUID) (bool, error) {
	fnExists := func() (bool, error) {
		return r.cacheByLogoID.Exists(logoID), nil
	}
	return r.Repository.Exists(fnExists)
}

// FindBySiteURL retrieves a site logo by its URL from the repository.
func (r *SiteLogoRepository) FindBySiteURL(_ context.Context, siteURL string) (*dmedia.SiteLogo, error) {
	// Define a function to find the site logo in the cache.
	fnFind := func() (*dmedia.SiteLogo, error) {
		return r.cacheBySiteURL.Find(siteURL), nil
	}
	return r.Repository.Find(fnFind)
}

// Checks if a site logo exists in the repository based on the provided site URL.
func (r *SiteLogoRepository) ExistsBySiteURL(_ context.Context, siteURL string) (bool, error) {
	fnExists := func() (bool, error) {
		return r.cacheBySiteURL.Exists(siteURL), nil
	}
	return r.Repository.Exists(fnExists)
}

// CleanExpired cleans expired items from the cache and calls the provided clean function.
func (r *SiteLogoRepository) CleanExpired(_ context.Context) error {
	clean := func() error {
		r.cacheByLogoID.CleanExpired()
		r.cacheBySiteURL.CleanExpired()
		return nil
	}
	return r.Repository.CleanExpired(clean)
}

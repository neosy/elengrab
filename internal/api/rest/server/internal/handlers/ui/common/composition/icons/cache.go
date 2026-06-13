package icons

import (
	"html/template"
	"time"

	memsimple "github.com/neosy/elengrab/internal/pkg/cache/memory/simple"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

const iconCacheTTL = 0 * time.Hour

var iconRep = newIconRepository(iconCacheTTL)

type iconEntry struct {
	raw template.HTML
}

func (icon *iconEntry) Copy() *iconEntry {
	return uptr.Copy(icon)
}

type iconRepository struct {
	memsimple.Repository[iconEntry]
	cache memsimple.Cache[string, iconEntry]
}

// newIconRepository returns a new object for the repository
func newIconRepository(ttl time.Duration) *iconRepository {
	r := &iconRepository{
		cache: memsimple.NewCacheWithDeaultCopier[string, iconEntry, *iconEntry](),
	}
	r.Repository.Init(ttl)
	return r
}

func (r *iconRepository) Save(fileName string, icon *iconEntry) error {
	if icon == nil || fileName == "" {
		return nil
	}

	save := func() error {
		r.cache.Save(
			fileName,
			icon,
			r.TTL(),
		)
		return nil
	}

	return r.Repository.Save(save)
}

func (r *iconRepository) SaveNegative(fileName string) error {
	if fileName == "" {
		return nil
	}

	save := func() error {
		r.cache.Save(
			fileName,
			nil,
			r.TTL(),
		)
		return nil
	}

	return r.Repository.Save(save)
}

// Delete removes a thumbnail from the in-memory repository using its ID.
func (r *iconRepository) Delete(fileName string) error {
	delete := func() error {
		if fileName != "" {
			r.cache.Delete(fileName)
		}
		return nil
	}
	return r.Repository.Delete(delete)
}

// FindByThumbID retrieves a thumbnail by its thumbnail ID from the repository.
func (r *iconRepository) Find(fileName string) (*iconEntry, memsimple.CacheStatus, error) {
	find := func() (*iconEntry, memsimple.CacheStatus, error) {
		data, status := r.cache.FindWithStatus(fileName)
		return data, status, nil
	}
	return r.Repository.FindWithStatus(find)
}

// Checks if a thumbnail exists by its thumbnail ID.
func (r *iconRepository) Exists(fileName string) (bool, error) {
	exists := func() (bool, error) {
		return r.cache.Exists(fileName), nil
	}
	return r.Repository.Exists(exists)
}

// CleanExpired cleans expired entries from the repository.
func (r *iconRepository) CleanExpired() error {
	// Define a clean function to remove expired entries from the cache.
	clean := func() error {
		r.cache.CleanExpired()
		return nil
	}
	// Call the base repository's CleanExpired method with the custom clean function.
	return r.Repository.CleanExpired(clean)
}

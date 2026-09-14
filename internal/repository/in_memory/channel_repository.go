package inmemory

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	memsimple "github.com/neosy/elengrab/internal/pkg/cache/memory/simple"
)

// Defines the structure for the in-memory repository of YouTube channels.
type ChannelRepository struct {
	// Embeds the base Repository
	memsimple.Repository[dmedia.Channel]

	// Cache for storing YouTube channels by their channel ID.
	cacheByChannelID         memsimple.Cache[uuid.UUID, dmedia.Channel]
	cacheByExternalChannelID memsimple.Cache[dtypes.ExternalChannelKey, dmedia.Channel]
}

// newChannelRepository returns a new object for the repository
func newChannelRepository(ttl time.Duration) *ChannelRepository {
	r := &ChannelRepository{
		cacheByChannelID:         memsimple.NewCacheWithDeaultCloner[uuid.UUID, dmedia.Channel, *dmedia.Channel](),
		cacheByExternalChannelID: memsimple.NewCacheWithDeaultCloner[dtypes.ExternalChannelKey, dmedia.Channel, *dmedia.Channel](),
	}
	r.Repository.Init(ttl)
	return r
}

func (r *ChannelRepository) Name() string {
	return "youtube_channel"
}

func (r *ChannelRepository) Save(ctx context.Context, channel *dmedia.Channel) error {
	if channel == nil {
		return errors.New("channel is nil")
	}

	if channel.ChannelID == uuid.Nil {
		return errors.New("channel ID is empty")
	}
	save := func() error {
		r.cacheByChannelID.Save(
			channel.ChannelID,
			channel,
			r.TTL(),
		)
		r.cacheByExternalChannelID.Save(
			dtypes.NewExternalChannelKey(channel.ExternalID, channel.Platform),
			channel,
			r.TTL(),
		)
		return nil
	}

	return r.Repository.Save(ctx, save)
}

func (r *ChannelRepository) SaveNegative(ctx context.Context, channelID uuid.UUID) error {
	if channelID == uuid.Nil {
		return nil
	}

	save := func() error {
		r.cacheByChannelID.Save(
			channelID,
			nil,
			r.TTL(),
		)
		return nil
	}

	return r.Repository.Save(ctx, save)
}

func (r *ChannelRepository) SaveNegativeByExternalChannelID(
	ctx context.Context,
	channelID string,
	platform string,
) error {
	if channelID == "" || platform == "" {
		return nil
	}

	save := func() error {
		r.cacheByExternalChannelID.Save(
			dtypes.NewExternalChannelKey(channelID, platform),
			nil,
			r.TTL(),
		)
		return nil
	}

	return r.Repository.Save(ctx, save)
}

// Delete removes a YouTube channel from the in-memory repository using its ID.
func (r *ChannelRepository) Delete(ctx context.Context, channelID uuid.UUID) error {
	delete := func() error {
		if channelID == uuid.Nil {
			return nil
		}

		channel := r.cacheByChannelID.Find(channelID)
		if channel != nil {
			r.cacheByExternalChannelID.Delete(channel.ExternalChannelKey())
			r.cacheByChannelID.Delete(channelID)
		}

		return nil
	}
	return r.Repository.Delete(ctx, delete)
}

// FindByChannelID retrieves a YouTube channel by its channel ID from the repository.
func (r *ChannelRepository) FindByChannelID(
	ctx context.Context,
	channelID uuid.UUID,
) (*dmedia.Channel, memsimple.CacheStatus, error) {
	find := func() (*dmedia.Channel, memsimple.CacheStatus, error) {
		data, status := r.cacheByChannelID.FindWithStatus(channelID)
		return data, status, nil
	}
	return r.Repository.FindWithStatus(ctx, find)
}

// Checks if a YouTube channel exists by its channel ID.
func (r *ChannelRepository) ExistsByChannelID(ctx context.Context, channelID uuid.UUID) (bool, error) {
	exists := func() (bool, error) {
		return r.cacheByChannelID.Exists(channelID), nil
	}
	return r.Repository.Exists(ctx, exists)
}

// CleanExpired cleans expired entries from the repository.
func (r *ChannelRepository) CleanExpired(ctx context.Context) error {
	// Define a clean function to remove expired entries from the cache.
	clean := func() error {
		r.cacheByChannelID.CleanExpired()
		r.cacheByExternalChannelID.CleanExpired()
		return nil
	}
	// Call the base repository's CleanExpired method with the custom clean function.
	return r.Repository.CleanExpired(ctx, clean)
}

func (r *ChannelRepository) FindByExternalChannelID(
	ctx context.Context,
	channelID string,
	platform string,
) (*dmedia.Channel, memsimple.CacheStatus, error) {
	find := func() (*dmedia.Channel, memsimple.CacheStatus, error) {
		data, status := r.cacheByExternalChannelID.FindWithStatus(dtypes.NewExternalChannelKey(channelID, platform))
		return data, status, nil
	}
	return r.Repository.FindWithStatus(ctx, find)
}

func (r *ChannelRepository) ExistsByExternalChannelID(
	ctx context.Context,
	channelID string,
	platform string,
) (bool, error) {
	exists := func() (bool, error) {
		return r.cacheByExternalChannelID.Exists(dtypes.NewExternalChannelKey(channelID, platform)), nil
	}
	return r.Repository.Exists(ctx, exists)
}

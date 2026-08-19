package inmemory

import (
	"context"
	"time"

	dmedia "github.com/neosy/elengrab/internal/domain/media"
	memsimple "github.com/neosy/elengrab/internal/pkg/cache/memory/simple"
)

// Defines the structure for the in-memory repository of YouTube channels.
type YoutubeChannelRepository struct {
	// Embeds the base Repository
	memsimple.Repository[dmedia.YoutubeChannel]

	// Cache for storing YouTube channels by their channel ID.
	cacheByChannel memsimple.Cache[string, dmedia.YoutubeChannel]
}

// newYoutubeChannelRepository returns a new object for the repository
func newYoutubeChannelRepository(ttl time.Duration) *YoutubeChannelRepository {
	r := &YoutubeChannelRepository{
		cacheByChannel: memsimple.NewCacheWithDeaultCopier[string, dmedia.YoutubeChannel, *dmedia.YoutubeChannel](),
	}
	r.Repository.Init(ttl)
	return r
}

func (r *YoutubeChannelRepository) Name() string {
	return "youtube_channel"
}

func (r *YoutubeChannelRepository) Save(ctx context.Context, channel *dmedia.YoutubeChannel) error {
	if channel == nil || channel.ChannelID == "" {
		return nil
	}

	save := func() error {
		r.cacheByChannel.Save(
			channel.ChannelID,
			channel,
			r.TTL(),
		)
		return nil
	}

	return r.Repository.Save(ctx, save)
}

func (r *YoutubeChannelRepository) SaveNegative(ctx context.Context, channelID string) error {
	if channelID == "" {
		return nil
	}

	save := func() error {
		r.cacheByChannel.Save(
			channelID,
			nil,
			r.TTL(),
		)
		return nil
	}

	return r.Repository.Save(ctx, save)
}

// Delete removes a YouTube channel from the in-memory repository using its ID.
func (r *YoutubeChannelRepository) Delete(ctx context.Context, channelID string) error {
	delete := func() error {
		if channelID != "" {
			r.cacheByChannel.Delete(channelID)
		}
		return nil
	}
	return r.Repository.Delete(ctx, delete)
}

// FindByChannelID retrieves a YouTube channel by its channel ID from the repository.
func (r *YoutubeChannelRepository) FindByChannelID(
	ctx context.Context,
	channelID string,
) (*dmedia.YoutubeChannel, memsimple.CacheStatus, error) {
	find := func() (*dmedia.YoutubeChannel, memsimple.CacheStatus, error) {
		data, status := r.cacheByChannel.FindWithStatus(channelID)
		return data, status, nil
	}
	return r.Repository.FindWithStatus(ctx, find)
}

// Checks if a YouTube channel exists by its channel ID.
func (r *YoutubeChannelRepository) ExistsByChannelID(ctx context.Context, channelID string) (bool, error) {
	exists := func() (bool, error) {
		return r.cacheByChannel.Exists(channelID), nil
	}
	return r.Repository.Exists(ctx, exists)
}

// CleanExpired cleans expired entries from the repository.
func (r *YoutubeChannelRepository) CleanExpired(ctx context.Context) error {
	// Define a clean function to remove expired entries from the cache.
	clean := func() error {
		r.cacheByChannel.CleanExpired()
		return nil
	}
	// Call the base repository's CleanExpired method with the custom clean function.
	return r.Repository.CleanExpired(ctx, clean)
}

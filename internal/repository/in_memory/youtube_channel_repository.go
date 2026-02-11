package inmemory

import (
	"context"
	"time"

	dmedia "github.com/neosy/elengrab/internal/domain/media"
	nmemory "github.com/neosy/elengrab/pkg/ncache/memory"
)

// Defines the structure for the in-memory repository of YouTube channels.
type YoutubeChannelRepository struct {
	// Embeds the base Repository
	nmemory.Repository[dmedia.YoutubeChannel]

	// Cache for storing YouTube channels by their channel ID.
	cacheByChannel nmemory.Cache[string, dmedia.YoutubeChannel]
}

// newYoutubeChannelRepository returns a new object for the repository
func newYoutubeChannelRepository(ttl time.Duration) *YoutubeChannelRepository {
	r := &YoutubeChannelRepository{
		cacheByChannel: nmemory.NewCacheWithDeaultCopier[string, dmedia.YoutubeChannel, *dmedia.YoutubeChannel](),
	}
	r.Repository.Init(ttl)
	return r
}

func (r *YoutubeChannelRepository) Save(_ context.Context, channel *dmedia.YoutubeChannel) error {
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

	return r.Repository.Save(save)
}

// Inserts a new YouTube channel into the in-memory repository.
func (r *YoutubeChannelRepository) Insert(ctx context.Context, channel *dmedia.YoutubeChannel) error {
	return r.Save(ctx, channel)
}

// Updates an existing YouTube channel into the in-memory repository.
func (r *YoutubeChannelRepository) Update(ctx context.Context, channel *dmedia.YoutubeChannel) error {
	return r.Save(ctx, channel)
}

// Delete removes a YouTube channel from the in-memory repository using its ID.
func (r *YoutubeChannelRepository) Delete(_ context.Context, channelID string) error {
	delete := func() error {
		r.cacheByChannel.Delete(channelID)
		return nil
	}
	return r.Repository.Delete(delete)
}

// FindByChannelID retrieves a YouTube channel by its channel ID from the repository.
func (r *YoutubeChannelRepository) FindByChannelID(_ context.Context, channelID string) (*dmedia.YoutubeChannel, error) {
	find := func() (*dmedia.YoutubeChannel, error) {
		return r.cacheByChannel.Find(channelID), nil
	}
	return r.Repository.Find(find)
}

// Checks if a YouTube channel exists by its channel ID.
func (r *YoutubeChannelRepository) ExistsByChannelID(_ context.Context, channelID string) (bool, error) {
	exists := func() (bool, error) {
		return r.cacheByChannel.Exists(channelID), nil
	}
	return r.Repository.Exists(exists)
}

// CleanExpired cleans expired entries from the repository.
func (r *YoutubeChannelRepository) CleanExpired(_ context.Context) error {
	// Define a clean function to remove expired entries from the cache.
	clean := func() error {
		r.cacheByChannel.CleanExpired()
		return nil
	}
	// Call the base repository's CleanExpired method with the custom clean function.
	return r.Repository.CleanExpired(clean)
}

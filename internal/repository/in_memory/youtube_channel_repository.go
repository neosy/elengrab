package inmemory

import (
	"context"
	"time"

	dmedia "github.com/neosy/elengrab/internal/domain/media"
	nmemory "github.com/neosy/elengrab/pkg/ncache/memory"
)

type YoutubeChannelRepository struct {
	nmemory.Repository[dmedia.YoutubeChannel]

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

func (r *YoutubeChannelRepository) Insert(ctx context.Context, channel *dmedia.YoutubeChannel) error {
	return r.Save(ctx, channel)
}

func (r *YoutubeChannelRepository) Update(ctx context.Context, channel *dmedia.YoutubeChannel) error {
	return r.Save(ctx, channel)
}

func (r *YoutubeChannelRepository) Delete(_ context.Context, channelID string) error {
	delete := func() error {
		r.cacheByChannel.Delete(channelID)
		return nil
	}
	return r.Repository.Delete(delete)
}

func (r *YoutubeChannelRepository) FindByChannelID(_ context.Context, channelID string) (*dmedia.YoutubeChannel, error) {
	find := func() (*dmedia.YoutubeChannel, error) {
		return r.cacheByChannel.Find(channelID), nil
	}
	return r.Repository.Find(find)
}

func (r *YoutubeChannelRepository) ExistsByChannelID(_ context.Context, channelID string) (bool, error) {
	exists := func() (bool, error) {
		return r.cacheByChannel.Exists(channelID), nil
	}
	return r.Repository.Exists(exists)
}

func (r *YoutubeChannelRepository) CleanExpired(_ context.Context) error {
	clean := func() error {
		r.cacheByChannel.CleanExpired()
		return nil
	}
	return r.Repository.CleanExpired(clean)
}

package inmemory

import (
	"context"
	"time"

	dyoutube "github.com/neosy/elengrab/internal/domain/youtube_info"
	nmemory "github.com/neosy/elengrab/pkg/ncache/memory"
	uptr "github.com/neosy/elengrab/pkg/utils/pointer"
)

type YoutubeChannelRepository struct {
	nmemory.Repository[dyoutube.YoutubeChannel]

	dataByChannelMap nmemory.Cache[string, dyoutube.YoutubeChannel]
}

// newYoutubeChannelRepository returns a new object for the repository
func newYoutubeChannelRepository(ttl time.Duration) *YoutubeChannelRepository {
	r := &YoutubeChannelRepository{
		dataByChannelMap: make(nmemory.Cache[string, dyoutube.YoutubeChannel]),
	}
	r.Repository.Init(ttl)
	return r
}

func (r *YoutubeChannelRepository) Save(_ context.Context, channel *dyoutube.YoutubeChannel) error {
	if channel == nil || channel.ChannelID == "" {
		return nil
	}

	save := func() error {
		r.dataByChannelMap.Save(
			channel.ChannelID,
			channel,
			r.copyChannel,
			r.TTL(),
		)
		return nil
	}

	return r.Repository.Save(save)
}

func (r *YoutubeChannelRepository) Insert(ctx context.Context, channel *dyoutube.YoutubeChannel) error {
	return r.Save(ctx, channel)
}

func (r *YoutubeChannelRepository) Update(ctx context.Context, channel *dyoutube.YoutubeChannel) error {
	return r.Save(ctx, channel)
}

func (r *YoutubeChannelRepository) Delete(_ context.Context, channelID string) error {
	delete := func() error {
		r.dataByChannelMap.Delete(channelID)
		return nil
	}
	return r.Repository.Delete(delete)
}

func (r *YoutubeChannelRepository) FindByChannelID(_ context.Context, channelID string) (*dyoutube.YoutubeChannel, error) {
	find := func() (*dyoutube.YoutubeChannel, error) {
		return r.dataByChannelMap.Find(channelID, r.copyChannel), nil
	}
	return find()
}

func (r *YoutubeChannelRepository) ExistsByChannelID(_ context.Context, channelID string) (bool, error) {
	exists := func() (bool, error) {
		return r.dataByChannelMap.Exists(channelID), nil
	}
	return exists()
}

func (r *YoutubeChannelRepository) copyChannel(channel *dyoutube.YoutubeChannel) *dyoutube.YoutubeChannel {
	if channel == nil {
		return nil
	}

	copy := uptr.Copy(channel)

	if len(channel.ImageRaw) > 1 {
		copy.ImageRaw = append([]byte{}, channel.ImageRaw...)
	}

	return copy
}

func (r *YoutubeChannelRepository) CleanExpired(_ context.Context) error {
	clean := func() error {
		r.dataByChannelMap.CleanExpired()
		return nil
	}
	return clean()
}

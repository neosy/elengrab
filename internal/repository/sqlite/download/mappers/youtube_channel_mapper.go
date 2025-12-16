package mappers

import (
	dyoutube "github.com/neosy/elengrab/internal/domain/youtube_info"
	edownload "github.com/neosy/elengrab/internal/repository/sqlite/download/entity"
)

func (m *Mappers) MapYoutubeChannelDomainToEntity(channel *dyoutube.YoutubeChannel) (*edownload.YoutubeChannel, error) {
	return &edownload.YoutubeChannel{
		ChannelID:   channel.ChannelID,
		ImageURL:    channel.ImageURL,
		ImageRaw:    channel.ImageRaw,
		ImageFormat: channel.ImageFormat,
	}, nil
}

func (m *Mappers) MapYoutubeChannelEntityToDomain(eChannel *edownload.YoutubeChannel) (*dyoutube.YoutubeChannel, error) {
	return &dyoutube.YoutubeChannel{
		ChannelID:   eChannel.ChannelID,
		ImageURL:    eChannel.ImageURL,
		ImageRaw:    eChannel.ImageRaw,
		ImageFormat: eChannel.ImageFormat,
		CreatedAt:   eChannel.CreatedAt,
		UpdatedAt:   eChannel.UpdatedAt,
	}, nil
}

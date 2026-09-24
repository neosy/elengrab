package deferred

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
)

func (m *migrations) fill_channel_info(ctx context.Context) (bool, error) {
	channelIDs := make(map[uuid.UUID]struct{})
	downloads := make([]*ddownload.MediaDownload, 0)

	err := m.Usecases().MediaDownload.IterateAll(ctx, func(d *ddownload.MediaDownload) error {
		if d.ChannelID == uuid.Nil {
			downloads = append(downloads, d)
			return nil
		}

		_, exists := channelIDs[d.ChannelID]
		if exists {
			return nil
		}

		channelIDs[d.ChannelID] = struct{}{}
		downloads = append(downloads, d)

		return nil
	})
	if err != nil {
		return false, err
	}

	if len(downloads) == 0 {
		return true, nil
	}

	updatedChannelIDs := make(map[uuid.UUID]struct{})

	for _, download := range downloads {
		channelSource, err := m.Services().Downloader.FetchChannelInfoWithCookieFallback(ctx, download.MediaURL)
		if err != nil {
			continue
		}

		if channelSource == nil || channelSource.ChannelID == "" || channelSource.Platform == "" {
			continue
		}

		channelID := download.ChannelID

		if channelID == uuid.Nil {
			channel, err := m.Usecases().Channel.FindByExternalChannelIDNoCache(ctx, channelSource.ChannelID, channelSource.Platform)
			if err != nil {
				return false, err
			}

			if channel != nil {
				channelID = channel.ChannelID
			}
		}

		if channelID == uuid.Nil {
			newChannel := dmedia.NewChannelFromSource(channelSource)

			if !newChannel.IsValid() {
				continue
			}

			err := m.Usecases().Channel.Create(ctx, newChannel)
			if err != nil {
				return false, err
			}

			channelID = newChannel.ChannelID
		} else if _, exists := updatedChannelIDs[channelID]; !exists {
			err = m.Usecases().Channel.Patch(ctx, channelID, func(channel *dmedia.Channel) error {
				if channel == nil {
					return nil
				}

				if channelSource.ChannelID != "" {
					channel.ExternalID = channelSource.ChannelID
				}

				if channelSource.Platform != "" {
					channel.Platform = channelSource.Platform
				}

				if channelSource.URL != "" {
					channel.ChannelURL = channelSource.URL
				}

				if channelSource.Host != "" {
					channel.Host = channelSource.Host
				}

				if channelSource.Username != "" {
					channel.Username = channelSource.Username
				}

				if channelSource.UsernameURL != "" {
					channel.UsernameURL = channelSource.UsernameURL
				}

				if channelSource.Title != "" {
					channel.Title = channelSource.Title
				}

				return nil
			})
			if err != nil {
				return false, err
			}
		}

		updatedChannelIDs[channelID] = struct{}{}

		if download.ChannelID == uuid.Nil {
			err = m.Usecases().MediaDownload.Patch(
				ctx, nil, download.DownloadID,
				func(d *ddownload.MediaDownload) error {
					if d == nil {
						return nil
					}

					d.ChannelID = channelID

					return nil
				},
			)
			if err != nil {
				return false, err
			}
		}

	}

	return true, nil
}

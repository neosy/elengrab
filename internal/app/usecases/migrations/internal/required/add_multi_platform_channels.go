package required

import (
	"context"
	"errors"

	"github.com/google/uuid"

	hostdetect "github.com/neosy/elengrab/internal/app/utils/host_detect"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (m *migrations) addMultiPlatformChannels(ctx context.Context) (bool, error) {
	m.Logger().Info("Adding support for multiple channel platforms...")

	var channels []*dmedia.Channel

	m.Usecases().Channel.IterateAll(
		ctx,
		func(channel *dmedia.Channel) error {
			channels = append(channels, channel)
			return nil
		},
	)

	if len(channels) == 0 {
		m.Logger().Info("No channels found to update")
		return true, nil
	}

	err := m.Usecases().Channel.Tx(ctx, func(ctx context.Context) error {
		channelIDsByExternalID := make(map[string]uuid.UUID)

		for _, channel := range channels {
			platformType := channel.PlatformType()
			newPlatformType := hostdetect.DetectPlatformType(channel.ChannelURL)

			if newPlatformType == dtypes.MediaPlatformTypeNone {
				newPlatformType = platformType
			}

			if newPlatformType == dtypes.MediaPlatformTypeNone {
				continue
			}

			oldID := channel.ChannelID
			newID := uuid.New()

			err := m.Usecases().Channel.UpdateChannelID(ctx, oldID, newID)
			if err != nil {
				return err
			}

			if newPlatformType != platformType {
				err = m.Usecases().Channel.Patch(ctx, newID, func(c *dmedia.Channel) error {
					if c == nil {
						m.Logger().Warn(
							"Channel not found",
							"channel_id", newID,
							"externalID", channel.ExternalID,
							"platform", channel.Platform,
						)
						return errors.New("channel not found")
					}

					c.Platform = newPlatformType.String()

					return nil
				})
			}

			channelIDsByExternalID[channel.ExternalID] = newID
		}

		err := m.Usecases().MediaDownload.Tx(ctx, func(ctx context.Context) error {
			for extID, newID := range channelIDsByExternalID {
				err := m.Usecases().MediaDownload.UpdateChannelID(ctx, extID, newID)
				if err != nil {
					return err
				}
			}

			return nil
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return false, err
	}

	m.Logger().Info("Added support for multiple channel platforms")

	return true, nil
}

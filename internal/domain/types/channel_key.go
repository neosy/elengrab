package dtypes

type ExternalChannelKey struct {
	// External channel identifier
	ChannelID string

	// Platform identifier
	Platform string
}

func NewExternalChannelKey(channelID, platform string) ExternalChannelKey {
	return ExternalChannelKey{
		ChannelID: channelID,
		Platform:  platform,
	}
}

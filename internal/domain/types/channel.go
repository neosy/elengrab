package dtypes

type ChannelSource struct {
	URL string

	ChannelID string
	Platform  string

	Host string

	Title string

	Image *ChannelImage
}

func (c *ChannelSource) Clone() *ChannelSource {
	if c == nil {
		return nil
	}

	channel := *c
	channel.Image = c.Image.Clone()

	return &channel
}

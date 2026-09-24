package dtypes

type ChannelSource struct {
	ChannelID string
	Platform  string

	URL  string
	Host string

	Username    string
	UsernameURL string

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

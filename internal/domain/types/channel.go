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

func (c *ChannelSource) Equal(other *ChannelSource) bool {
	if c == nil || other == nil {
		return c == other
	}

	return c.ChannelID == other.ChannelID &&
		c.Platform == other.Platform &&
		c.URL == other.URL &&
		c.Host == other.Host &&
		c.Username == other.Username &&
		c.UsernameURL == other.UsernameURL &&
		c.Title == other.Title &&
		c.Image.Equal(other.Image)
}

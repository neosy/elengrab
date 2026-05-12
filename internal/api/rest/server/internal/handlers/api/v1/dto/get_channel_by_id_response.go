package dto

type GetChannelByIDResponse struct {
	ChannelID   string `json:"channelId"`
	ImageURL    string `json:"imageUrl,omitempty"`
	ImageFormat string `json:"imageFormat,omitempty"`
}

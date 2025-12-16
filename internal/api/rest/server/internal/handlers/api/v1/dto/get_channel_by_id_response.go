package dto

type GetChannelByIDResponse struct {
	ChannelID   string `json:"channelID"`
	ImageURL    string `json:"imageURL,omitempty"`
	ImageFormat string `json:"imageFormat,omitempty"`
}

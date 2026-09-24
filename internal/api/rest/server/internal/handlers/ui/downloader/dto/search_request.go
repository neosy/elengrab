package dto

type SearchRequest struct {
	QueryText string `json:"query" validate:"max=100"`
	ViewMode  string `json:"viewMode" validate:"omitempty,viewMode"`

	ChannelID string `json:"channelId" validate:"omitempty"`
}

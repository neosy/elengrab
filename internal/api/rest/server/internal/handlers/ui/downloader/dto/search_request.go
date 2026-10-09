package dto

type SearchRequest struct {
	QueryText string `json:"query" validate:"max=100"`
	ViewMode  string `json:"viewMode" validate:"omitempty,viewMode"`

	UserName  string `json:"userName" validate:"omitempty"`
	ChannelID string `json:"channelId" validate:"omitempty"`
}

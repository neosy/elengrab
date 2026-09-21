package dto

type SearchRequest struct {
	Query    string `json:"query" validate:"max=100"`
	ViewMode string `json:"viewMode" validate:"omitempty,viewMode"`

	ChannelID string `json:"channelId" validate:"omitempty"`
}

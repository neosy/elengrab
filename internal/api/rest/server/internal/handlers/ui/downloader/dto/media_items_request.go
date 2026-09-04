package dto

type MediaItemsRequest struct {
	ViewMode string `json:"viewMode" validate:"required,viewMode"`
	Search   string `json:"search" validate:"max=100"`

	LastID        string `json:"lastId" validate:"omitempty"`
	LastCreatedAt string `json:"lastCreatedAt" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	LastViews     int    `json:"lastViews" validate:"omitempty,gte=0"`
}

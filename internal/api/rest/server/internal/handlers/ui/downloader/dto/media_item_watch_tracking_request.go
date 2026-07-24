package dto

type MediaItemWatchTrackingRequest struct {
	PositionMs int `json:"positionMs" validate:"required,gt=0"`
	IntervalMs int `json:"intervalMs" validate:"required,gt=0"`
}

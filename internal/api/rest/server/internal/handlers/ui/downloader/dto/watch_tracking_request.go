package dto

type MediaItemWatchTrackingRequest struct {
	PositionMs int    `json:"positionMs" validate:"omitempty,gte=0"`
	IntervalMs int    `json:"intervalMs" validate:"required,gt=0"`
	EventType  string `json:"eventType" validate:"omitempty,oneof=pause ended seek heartbeat"`
}

package dto

type GetLastWatchPositionResponse struct {
	DownloadID string `json:"itemId"`
	Position   uint32 `json:"position"`
}

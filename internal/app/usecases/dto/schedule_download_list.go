package dto

import (
	"uuid"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

type ScheduleDownloadListRequest struct {
	Items   []ScheduleDownloadListRequestItem
	Options ddownload.DownloadOptions
}

type ScheduleDownloadListResponse struct {
	Items []ScheduleDownloadListResponseItem
}

type ScheduleDownloadListResponseItem struct {
	DownloadID uuid.UUID

	URL        string
	MediaTitle string
}

type ScheduleDownloadListRequestItem struct {
	URL   string
	Title string
}

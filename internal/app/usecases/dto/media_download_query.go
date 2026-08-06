package dto

import "time"

type MediaDownloadFilters struct {
	Title string
}

type MediaDownloadQuery struct {
	Before time.Time
	Limit  uint64

	Filters MediaDownloadFilters
}

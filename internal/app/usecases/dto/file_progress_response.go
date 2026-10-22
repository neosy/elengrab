package dto

import (
	"uuid"
)

type MediaDownloadProgressResponse struct {
	DownloadID uuid.UUID
	Percent    int
}

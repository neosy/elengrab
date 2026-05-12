package dto

import (
	"github.com/google/uuid"
)

type MediaDownloadProgressResponse struct {
	DownloadID uuid.UUID
	Percent    int
}

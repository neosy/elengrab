package dto

import (
	"github.com/google/uuid"
)

type FileProgressResponse struct {
	FileID  uuid.UUID
	Percent int
}

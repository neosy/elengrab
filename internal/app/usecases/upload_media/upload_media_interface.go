package uploadmedia

import (
	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/upload_media/internal/upload"
)

type UploadMedia interface {
	NewUpload(userID uuid.UUID, originalName string) *upload.Upload
}

package uploadmedia

import (
	"uuid"

	"github.com/neosy/elengrab/internal/app/usecases/upload_media/internal/upload"
	uploadfile "github.com/neosy/elengrab/internal/app/usecases/upload_media/internal/upload"
)

func (u *uploadMedia) NewUpload(userID uuid.UUID, originalName string) *upload.Upload {
	return uploadfile.NewUpload(u.logger, userID, originalName, u.uploadDir, u.ffmpegService)
}

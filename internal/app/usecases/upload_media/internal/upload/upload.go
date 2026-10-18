package upload

import (
	"log/slog"

	"github.com/google/uuid"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
	pservices "github.com/neosy/elengrab/internal/ports/services"
)

type Upload struct {
	logger *slog.Logger

	ffmpegService pservices.FFMpeg

	uploadDir        string
	originalFileName string

	userID uuid.UUID

	fileID   uuid.UUID
	fileName string
	fileExt  string

	mediaInfo *dtypes.MediaInfo
}

func NewUpload(
	logger *slog.Logger,
	userID uuid.UUID,
	originalName, uploadDir string,
	ffmpegService pservices.FFMpeg,
) *Upload {
	fileID := newFileID()

	return &Upload{
		logger: logger,

		ffmpegService: ffmpegService,

		uploadDir:        uploadDir,
		originalFileName: originalName,

		userID: userID,

		fileID:   fileID,
		fileName: fileID.String(),
	}
}

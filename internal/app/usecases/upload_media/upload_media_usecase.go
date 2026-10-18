package uploadmedia

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/neosy/elengrab/internal/app/usecases/upload_media/internal/consts"
	pservices "github.com/neosy/elengrab/internal/ports/services"
)

type uploadMedia struct {
	logger *slog.Logger

	ffmpegService pservices.FFMpeg

	uploadDir string
}

func NewUploadMedia(logger *slog.Logger, storageDir string, deps Dependencies) UploadMedia {
	uploadDir := filepath.Join(storageDir, consts.TmpDirName)

	// Best effort at startup: a problem shows up in the log right away.
	// The directory is created again before every write, so a failure here is not fatal.
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		logger.Error("Failed to create the upload directory", "uploadDir", uploadDir, "error", err)
	}

	return &uploadMedia{
		logger: logger,

		ffmpegService: deps.FFmpegService,

		uploadDir: uploadDir,
	}
}

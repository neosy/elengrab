package upload

import (
	"path/filepath"

	"github.com/neosy/elengrab/internal/app/usecases/upload_media/internal/consts"
)

func (u *Upload) PartPath() string {
	return filepath.Join(u.uploadDir, u.fileName) + consts.PartFileExt
}

// FilePath returns the full path of the uploaded file, including its extension.
func (u *Upload) FilePath() string {
	return filepath.Join(u.uploadDir, u.fileName) + u.fileExt
}

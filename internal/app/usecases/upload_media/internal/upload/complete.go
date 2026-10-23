package upload

import (
	"context"
	"os"

	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/valyala/fasthttp"
)

// finish renames the .part file to its final name in the upload directory.
func (u *Upload) finish() error {
	if u.mediaInfo == nil {
		return errorx.New("media info is not available",
			errorx.WithErrorMessage("Failed to read the file's metadata."))
	}

	u.fileExt = u.mediaInfo.Format.Ext()

	err := os.Rename(u.PartPath(), u.FilePath())
	if err != nil {
		u.logger.Warn(
			"Failed to rename part file to final file",
			"partPath", u.PartPath(),
			"finalPath", u.FilePath(),
			"error", err,
		)
		return errorx.Errorf("failed to rename part file to final file: %w", err,
			errorx.WithErrorMessage("Failed to save the file. Please try again."))
	}

	return nil
}

// Complete validates the uploaded .part file and renames it to its final name
// in the upload directory. On any error the file stays a .part file; call Discard to remove it.
func (u *Upload) Complete(ctx context.Context) error {
	mediaInfo, err := u.extractMediaInfo(ctx)
	if err != nil {
		return err
	}

	u.mediaInfo = mediaInfo

	if err := u.validate(); err != nil {
		return errorx.NewFromError(err, errorx.WithHttpStatus(fasthttp.StatusUnprocessableEntity))
	}

	err = u.finish()
	if err != nil {
		return errorx.NewFromError(err, errorx.WithHttpStatus(fasthttp.StatusInternalServerError))
	}

	return nil
}

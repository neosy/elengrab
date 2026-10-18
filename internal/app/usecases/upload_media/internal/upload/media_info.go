package upload

import (
	"context"
	"errors"

	apperrors "github.com/neosy/elengrab/internal/app/errors"
	"github.com/neosy/elengrab/internal/app/usecases/upload_media/internal/consts"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/valyala/fasthttp"
)

// MediaInfo returns the media information of the uploaded file.
func (u *Upload) MediaInfo() *dtypes.MediaInfo {
	return u.mediaInfo
}

func (u *Upload) extractMediaInfo(ctx context.Context) (*dtypes.MediaInfo, error) {
	filePath := u.PartPath()

	probeCtx, cancel := context.WithTimeout(ctx, consts.ProbeTimeout)
	defer cancel()

	mediaInfo, err := u.ffmpegService.ExtractMediaInfoWithFFprobe(probeCtx, filePath, nil)
	if err != nil {
		u.logger.Warn(
			"Failed to get media info from file",
			"filePath", filePath,
			"error", err,
		)

		// A timeout is a server-side problem, not a problem of the file.
		if errors.Is(probeCtx.Err(), context.DeadlineExceeded) {
			return nil, errorx.Errorf("ffprobe timed out: %w", err,
				errorx.WithErrorMessage("Reading the file took too long. Please try again."),
				errorx.WithHttpStatus(fasthttp.StatusInternalServerError),
			)
		}

		if errors.Is(err, apperrors.ErrUnsupportedVideoCodec) ||
			errors.Is(err, apperrors.ErrUnsupportedAudioCodec) ||
			errors.Is(err, apperrors.ErrUnsupportedMediaFormat) {
			return nil, err
		}

		// ffprobe could not parse the file: it is corrupted or not a media file.
		return nil, errorx.Errorf("failed to extract media info with ffprobe: %w", err,
			errorx.WithErrorMessage("Failed to read the file's metadata."),
			errorx.WithHttpStatus(fasthttp.StatusUnprocessableEntity),
		)
	}

	return mediaInfo.ToDomain(), nil
}

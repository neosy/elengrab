package upload

import (
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/errorx"
)

// validate checks the parsed media info.
func (u *Upload) validate() error {
	if u.mediaInfo == nil {
		return errorx.New("media info is not available",
			errorx.WithErrorMessage("Failed to read the file's metadata."))
	}

	if u.mediaInfo.VideoInfo == nil && u.mediaInfo.AudioInfo == nil {
		return errorx.New("file contains neither video nor audio streams",
			errorx.WithErrorMessage("The file contains no audio or video."))
	}

	if u.mediaInfo.VideoInfo != nil && u.mediaInfo.VideoInfo.Codec == dtypes.VideoCodecNone {
		return errorx.New("video stream has an unknown or unsupported codec",
			errorx.WithErrorMessage("The video format is not supported."))
	}

	if u.mediaInfo.AudioInfo != nil && u.mediaInfo.AudioInfo.Codec == dtypes.AudioCodecNone {
		return errorx.New("audio stream has an unknown or unsupported codec",
			errorx.WithErrorMessage("The audio format is not supported."))
	}

	if u.mediaInfo.Duration() <= 0 {
		return errorx.New("media has no duration",
			errorx.WithErrorMessage("Could not determine the media duration."))
	}

	if v := u.mediaInfo.VideoInfo; v != nil && (v.Width <= 0 || v.Height <= 0) {
		return errorx.New("video stream has invalid dimensions",
			errorx.WithErrorMessage("The video has invalid dimensions."))
	}

	return nil
}

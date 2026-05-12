package statussetter

import "log/slog"

type MediaDownloadStatusSetter struct {
	logger *slog.Logger
}

// NewMediaDownloadStatusSetter creates a new instance of MediaDownloadStatusSetter.
func NewMediaDownloadStatusSetter(logger *slog.Logger) *MediaDownloadStatusSetter {
	return &MediaDownloadStatusSetter{
		logger: logger,
	}

}

package statussetter

import "log/slog"

type DownloadTaskStatusSetter struct {
	logger *slog.Logger
}

// NewFDownloadTaskStatusSetter creates a new instance of DownloadTaskStatusSetter.
func NewFDownloadTaskStatusSetter(logger *slog.Logger) *DownloadTaskStatusSetter {
	return &DownloadTaskStatusSetter{
		logger: logger,
	}

}

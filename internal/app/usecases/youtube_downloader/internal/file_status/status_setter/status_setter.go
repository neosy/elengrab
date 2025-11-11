package statussetter

import "log/slog"

type FileStatusSetter struct {
	logger *slog.Logger
}

// NewFileStatusSetter creates a new instance of FileStatusSetter.
func NewFileStatusSetter(logger *slog.Logger) *FileStatusSetter {
	return &FileStatusSetter{
		logger: logger,
	}

}

package statussetter

import "log/slog"

type StatusSetter struct {
	logger *slog.Logger
}

// NewStatusSetter creates a new instance of StatusSetter.
func NewStatusSetter(logger *slog.Logger) *StatusSetter {
	return &StatusSetter{
		logger: logger,
	}
}

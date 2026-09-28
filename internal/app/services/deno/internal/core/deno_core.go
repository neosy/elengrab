package core

import "log/slog"

type DenoCore struct {
	logger *slog.Logger

	// parameters
	denoPath string
}

func NewDenoCore(logger *slog.Logger, denoPath string) *DenoCore {
	return &DenoCore{
		logger: logger,

		// parameters
		denoPath: denoPath,
	}
}

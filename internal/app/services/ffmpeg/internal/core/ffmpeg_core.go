package core

import "log/slog"

type FFmpegCore struct {
	logger *slog.Logger

	info *info

	// parameters
	ffmpegPath string
}

func NewFFmpegCore(logger *slog.Logger, path string) *FFmpegCore {
	return &FFmpegCore{
		logger: logger,

		info: newInfo(),

		// parameters
		ffmpegPath: path,
	}
}

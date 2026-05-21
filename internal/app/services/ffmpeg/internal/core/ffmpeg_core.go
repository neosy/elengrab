package core

import "log/slog"

type FFmpegCore struct {
	logger *slog.Logger

	info *info

	// parameters
	ffmpegPath  string
	ffprobePath string
}

func NewFFmpegCore(logger *slog.Logger, ffmpegPath, ffprobePath string) *FFmpegCore {
	return &FFmpegCore{
		logger: logger,

		info: newInfo(),

		// parameters
		ffmpegPath:  ffmpegPath,
		ffprobePath: ffprobePath,
	}
}

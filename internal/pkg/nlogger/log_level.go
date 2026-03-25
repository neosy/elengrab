package nlogger

import "log/slog"

type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

var mapLogLevelToSlogLevel = map[LogLevel]slog.Level{
	LogLevelDebug: slog.LevelDebug,
	LogLevelInfo:  slog.LevelInfo,
	LogLevelWarn:  slog.LevelWarn,
	LogLevelError: slog.LevelError,
}

func LevelToSlogLevel(level LogLevel) slog.Level {
	return mapLogLevelToSlogLevel[level]
}

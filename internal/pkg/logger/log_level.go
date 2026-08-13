package logger

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

// LevelToSlogLevel converts a LogLevel to the corresponding slog.Level.
func LevelToSlogLevel(level LogLevel) slog.Level {
	return mapLogLevelToSlogLevel[level]
}

// String returns the string representation of the LogLevel.
func (level LogLevel) String() string {
	return string(level)
}

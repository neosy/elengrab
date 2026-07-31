package nlogger

type LogFormat string

const (
	LogFormatJson    LogFormat = "json"
	LogFormatConsole LogFormat = "console"
)

// String returns the string representation of the LogFormat.
func (lf LogFormat) String() string {
	return string(lf)
}

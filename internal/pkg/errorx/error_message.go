package errorx

// ErrorMessageProvider returns an error message string.
type ErrorMessageProvider func() *string

// ErrorMessageArg creates a provider that always returns the given text.
func ErrorMessageArg(text string) ErrorMessageProvider {
	return func() *string {
		return &text
	}
}

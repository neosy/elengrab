package types

import "github.com/neosy/elengrab/internal/pkg/stringx"

// ErrorMessageProvider returns an error message string.
type ErrorMessageProvider func() *string

// ErrorMessageArg creates a provider that always returns the given text.
func ErrorMessageArg(text string) ErrorMessageProvider {
	return func() *string {
		msg := stringx.Capitalize(text)
		return &msg
	}
}

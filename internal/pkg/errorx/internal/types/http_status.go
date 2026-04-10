package types

type (
	// HttpStatusProvider is a function that provides an HTTP status.
	HttpStatusProvider func() *int
)

// WithHttpStatus returns a function that provides the HTTP status.
func WithHttpStatus(status int) HttpStatusProvider {
	return func() *int {
		return &status
	}
}

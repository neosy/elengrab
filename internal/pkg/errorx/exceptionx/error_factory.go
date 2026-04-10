package exceptionx

// ErrorxFactory defines the interface for creating errors from exceptions.
type ErrorxFactory interface {
	// New creates a new error with the given text and optional arguments.
	New(text string, args ...any) error
	// NewFromException creates a new error based on the provided Exception and optional arguments.
	NewFromException(ex Exception, args ...any) error
	// NewFromDomainException creates a new error based on the provided DomainException and optional arguments.
	NewFromDomainException(ex DomainException, args ...any) error
}

var errorxFactory ErrorxFactory

// RegisterFactory registers the provided ErrorxFactory for creating errors from exceptions.
func RegisterFactory(f ErrorxFactory) {
	errorxFactory = f
}

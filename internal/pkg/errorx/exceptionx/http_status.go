package exceptionx

// HttpStatusProvider is a function that provides an HTTP status code.
type HttpStatusProvider func() *int

// HttpStatusArg returns a function that provides the HTTP status code.
func HttpStatusArg(code int) HttpStatusProvider {
	return func() *int {
		return &code
	}
}

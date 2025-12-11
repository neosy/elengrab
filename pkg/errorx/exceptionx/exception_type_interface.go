package exceptionx

// Interface type exception type
type ExceptionType interface {
	String() string
	HttpStatusCode() int
}

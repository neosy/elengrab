package exceptionx

type HttpStatusProvider func() *int

func HttpStatusArg(code int) HttpStatusProvider {
	return func() *int {
		return &code
	}
}

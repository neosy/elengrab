package errorx

type HttpStatusProvider func() int

func ArgHttpStatusCode(code int) HttpStatusProvider {
	return func() int {
		return code
	}
}

package errorx

type errorType uint8

const (
	errorTypeMain errorType = iota
	errorTypeWrap
	errorTypeJoin
)

func typeOf(err error) errorType {
	switch err.(type) {
	case interface{ Unwrap() error }:
		return errorTypeWrap
	case interface{ Unwrap() []error }:
		return errorTypeJoin
	default:
		return errorTypeMain
	}
}

package exceptionx

// Interface type exception code
type ExceptionCode interface {
	Num() uint
	String() string
	Type() ExceptionType
}

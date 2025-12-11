package exceptionx

// Type Exception
type Exception struct {
	code  ExceptionCode
	eType ExceptionType
}

// NewException creating an Exception object
func NewException(
	code ExceptionCode,
	eType ExceptionType,
) *Exception {
	return &Exception{
		code:  code,
		eType: eType,
	}
}

// Code getting exception code
func (e *Exception) Code() (code ExceptionCode) {
	if e != nil {
		code = e.code
	}

	return
}

// Type getting exception type
func (e *Exception) Type() (eType ExceptionType) {
	if e != nil {
		eType = e.eType
	}

	return
}

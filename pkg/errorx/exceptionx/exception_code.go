package exceptionx

// Type exception code
type exceptionCode struct {
	num   uint
	text  string
	eType ExceptionType
}

// NewExceptionCode creating an exceptionCode object
func NewExceptionCode(num uint, text string, eType ExceptionType) *exceptionCode {
	return &exceptionCode{
		num:   num,
		text:  text,
		eType: eType,
	}
}

// Num getting the numeric value of the code of exception
func (e *exceptionCode) Num() (num uint) {
	if e != nil {
		num = e.num
	}

	return
}

// String getting text value of code
func (e *exceptionCode) String() (text string) {
	if e != nil {
		text = e.text
	}

	return
}

// Type getting ExceptionType for code of exception
func (e *exceptionCode) Type() (eType ExceptionType) {
	if e != nil {
		eType = e.eType
	}

	return
}

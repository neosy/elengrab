package exceptionx

// Type of exception codes
type exceptionCodes struct {
	count uint
	codes []ExceptionCode
}

// NewExceptionCodes creating an exceptionCodes object
func NewExceptionCodes(maxCount uint) exceptionCodes {
	exceptions := exceptionCodes{}
	exceptions.codes = make([]ExceptionCode, maxCount)

	return exceptions
}

// AddNum adding a new exception code by specifying a numeric code
func (e *exceptionCodes) AddByNum(num uint, text string, eType ExceptionType) ExceptionCode {
	e.count = num

	e.codes[e.count] = NewExceptionCode(e.count, text, eType)

	return e.codes[e.count]
}

// Add adding a new exception code
func (e *exceptionCodes) Add(text string, eType ExceptionType) ExceptionCode {
	e.count++

	e.codes[e.count] = NewExceptionCode(e.count, text, eType)

	return e.codes[e.count]
}

package exceptionx

// Type of exceptions
type exceptions struct {
	count     uint
	exception map[uint]Exception
}

// NewExceptions creating an exceptions object
func NewExceptions(maxCount uint) exceptions {
	exceptions := exceptions{}
	exceptions.exception = make(map[uint]Exception, maxCount)
	return exceptions
}

// AddNum adding a new exception by specifying a numeric code
func (e *exceptions) AddNum(num uint, code string, opts ...any) Exception {
	e.count = num
	e.exception[e.count] = NewException(e.count, code, opts...)
	return e.exception[e.count]
}

// Add adding a new exception sequentially
func (e *exceptions) AddSeq(code string, opts ...any) Exception {
	e.count++
	return e.AddNum(e.count, code, opts...)
}

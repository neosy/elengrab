package exceptionx

// Type of exceptions
type exceptions struct {
	count     uint
	exception []Exception
}

// NewExceptions creating an exceptions object
func NewExceptions(maxCount uint) exceptions {
	exceptions := exceptions{}
	exceptions.exception = make([]Exception, maxCount)
	return exceptions
}

// AddNum adding a new exception by specifying a numeric code
func (e *exceptions) AddNum(num uint, args ...any) Exception {
	e.count = num
	e.exception[e.count] = NewException(e.count, args)
	return e.exception[e.count]
}

// Add adding a new exception sequentially
func (e *exceptions) AddSeq(args ...any) Exception {
	e.count++
	return e.AddNum(e.count, args)
}

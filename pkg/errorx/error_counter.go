package errorx

// Type error counter
type errorCounter struct {
	count uint
}

// NewErrorCounter creating an errorCounter object
func NewErrorCounter(startNum uint) ErrorCounter {
	return &errorCounter{
		count: startNum,
	}
}

// Set sets and returns the counter number
func (c *errorCounter) Set(num uint) uint {
	c.count = num

	return c.count
}

// Inc increases the counter by one and returns the current value
func (c *errorCounter) Inc() uint {
	c.count++

	return c.count
}

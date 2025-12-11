package errorx

// Type error counter
type errorxCounter struct {
	count uint
}

// NewErrorxCounter creating an errorxCounter object
func NewErrorxCounter(startNum uint) ErrorxCounter {
	return &errorxCounter{
		count: startNum,
	}
}

// Set sets and returns the counter number
func (c *errorxCounter) Set(num uint) uint {
	c.count = num

	return c.count
}

// Inc increases the counter by one and returns the current value
func (c *errorxCounter) Inc() uint {
	c.count++

	return c.count
}

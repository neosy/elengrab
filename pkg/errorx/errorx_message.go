package errorx

// Type for error messages
type errorxMessage string

// NewErrorxMessage creating an errorxMessage object
func NewErrorxMessage(text string) ErrorxMessage {
	return (*errorxMessage)(&text)
}

// String returns text
func (m *errorxMessage) String() string {
	return string(*m)
}

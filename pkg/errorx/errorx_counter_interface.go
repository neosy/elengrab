package errorx

// Interface error counter
type ErrorxCounter interface {
	Set(num uint) uint
	Inc() uint
}

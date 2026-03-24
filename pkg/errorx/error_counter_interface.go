package errorx

// Interface error counter
type ErrorCounter interface {
	Set(num uint) uint
	Inc() uint
}

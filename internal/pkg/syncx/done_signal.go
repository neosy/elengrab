package syncx

import (
	"sync/atomic"
)

type DoneSignal struct {
	ch     chan struct{}
	closed atomic.Bool
}

// NewDoneSignal creates a done signal channel and a function to close it safely.
func NewDoneSignal() *DoneSignal {
	return &DoneSignal{
		ch: make(chan struct{}),
	}
}

// Done returns the done signal channel.
func (d *DoneSignal) Done() <-chan struct{} {
	return d.ch
}

// Close closes the done signal channel safely.
func (d *DoneSignal) Close() {
	if d.closed.CompareAndSwap(false, true) {
		close(d.ch)
	}
}

// IsClosed checks if the done signal channel is closed.
func (d *DoneSignal) IsClosed() bool {
	return d.closed.Load()
}

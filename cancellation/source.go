package cancellation

import "sync"
import "time"

// Source represents Cancellation source
type Source interface {
	// IsCancelled return true if Source already cancelled false otherwise
	IsCancelled() bool
	// Cancellation returns Cancellation from this Source
	Cancellation() Cancellation
	// Channel return channel which closes when Source cancelled
	Channel() <-chan struct{}
	// Cancel cancels Source and it's Cancellations.
	// Returns true if Source cancelled by this function call, false if Source was already cancelled
	Cancel() bool
}

type cancellationSource struct {
	ch chan struct{}
	m  sync.Mutex
}

func (cs *cancellationSource) IsCancelled() bool {
	select {
	case <-cs.ch:
		return true
	default:
		return false
	}
}

func (cs *cancellationSource) Channel() <-chan struct{} {
	return cs.ch
}

func (cs *cancellationSource) Cancellation() Cancellation {
	return cs
}

func (cs *cancellationSource) Cancel() bool {
	if cs.IsCancelled() {
		return false
	}

	cs.m.Lock()
	defer cs.m.Unlock()

	if cs.IsCancelled() {
		return false
	}

	close(cs.ch)
	return true
}

// NewSource creates new Source
func NewSource() Source {
	return &cancellationSource{
		ch: make(chan struct{}),
	}
}

// NewSourceWithTimeout creates new Source which cancelled after duration time
func NewSourceWithTimeout(duration time.Duration) Source {
	ret := NewSource()

	time.AfterFunc(duration, func() {
		ret.Cancel()
	})

	return ret
}

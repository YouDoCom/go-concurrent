package cancellation

// Cancellation represents async operation cancellation
type Cancellation interface {
	// IsCancelled return true if Cancellation cancelled false otherwise
	IsCancelled() bool
	// Channel return channel which closes when Cancellation cancelled
	Channel() <-chan struct{}
}

var cancelled, none Cancellation

// Cancelled returns cancelled Cancellation
func Cancelled() Cancellation {
	return cancelled
}

// None returns Cancellation which never will be cancelled
func None() Cancellation {
	return none
}

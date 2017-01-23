package cancellation

func init() {
	cs := NewSource()
	cs.Cancel()

	cancelled = cs.Cancellation()
	none = NewSource().Cancellation()
}

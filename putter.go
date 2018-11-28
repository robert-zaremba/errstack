package errstack

// StubPutter ignores all input. It only contains information whether any input was submitted
type StubPutter struct {
	hasError bool
}

// HasError checks if any error occurred
func (rp *StubPutter) HasError() bool {
	return rp.hasError
}

// Fork does nothing
func (rp *StubPutter) Fork(_ string) Putter {
	return rp
}

// ForkIdx does nothing
func (rp *StubPutter) ForkIdx(_ int) Putter {
	return rp
}

// Put adds new error
func (rp *StubPutter) Put(_ interface{}) {
	rp.hasError = true
}

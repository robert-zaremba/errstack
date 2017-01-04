package errstack

import (
	"errors"
	"net/url"
	"strings"
)

// UntilFirst is a struct to easily chain a sequence of operations
// until a first error arise
type UntilFirst struct {
	Err E
}

// Do performs an action
func (ue *UntilFirst) Do(f func() E) *UntilFirst {
	if ue.Err != nil {
		return ue
	}
	ue.Err = f()
	return ue
}

// UntilFirstDo is utility method for using UntilFirst
func UntilFirstDo(f func() E) *UntilFirst {
	uf := &UntilFirst{}
	return uf.Do(f)
}

// Seq returns a standard error from a string sequence
func Seq(ss ...string) error {
	return errors.New(strings.Join(ss, " "))
}

// IsTimeout checks if given error is a timeout (net/url.Error)
func IsTimeout(e error) bool {
	urlError, ok := e.(*url.Error)
	return ok && urlError.Timeout()
}

// Logger defines the reporting interface for utils error functions
type Logger interface {
	Error(msg string, ctx ...interface{})
}

// CallAndLog calls function which may return error and logs it.
// The intention of this function is to be used with `go` and `defer` clauses.
func CallAndLog(l Logger, f func() error) {
	Log(l, f())
}

// Log logs error if it's not nil
func Log(l Logger, err error) {
	if err != nil {
		l.Error("Unhandled error", err)
	}
}

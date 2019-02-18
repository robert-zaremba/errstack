package errstack

import (
	"encoding/json"
	"fmt"

	"github.com/facebookgo/stack"
)

// HasUnderlying describes entity (usually an error) which has underlying error.
type HasUnderlying interface {
	Cause() error
}

// HasStatusCode provides a function to return the HTTP status code.
type HasStatusCode interface {
	StatusCode() int
}

// HasStacktrace provides a function to return the the root stacktrace
type HasStacktrace interface {
	Stacktrace() stack.Stack
}

// E is error with more information. It is able to marshall itself to json as response.
// Result of Error() method should include stacktrace. Therefore it should not be
// displayed directly to the user
type E interface {
	error
	HasStatusCode
	HasStacktrace
	json.Marshaler
	IsReq() bool
	WithMsg(string) E
}

type errstack struct {
	err        error
	stacktrace stack.Stack
	message    string
}

func newErr(e error, s string, skip int) errstack {
	st := stack.Callers(skip + 1)
	return errstack{e, st, s}
}

func (e *errstack) Error() string {
	var message = e.message
	if message == "" {
		message = "error"
	}
	if e.err == nil {
		return message
	}
	return fmt.Sprint(message, " [", e.err.Error(), "]")
}

func (e errstack) withMsg(msg string) errstack {
	return errstack{
		err:        wrapper{e.message, e.err},
		message:    msg,
		stacktrace: e.stacktrace,
	}
}

// Cause implements HasCause interface
func (e *errstack) Cause() error {
	return e.err
}

// Stacktrace returns error creation stacktrace
func (e *errstack) Stacktrace() stack.Stack {
	return e.stacktrace
}

type wrapper struct {
	msg string
	err error
}

func (e wrapper) Error() string {
	errmsg := "nil"
	if e.err != nil {
		errmsg = e.err.Error()
	}
	return fmt.Sprintf("%s [%s]", e.msg, errmsg)
}

func (e *wrapper) Cause() error {
	return Cause(e.err)
}

// Cause returns the underlying cause of the error, if possible.
// An error value has a cause if it implements the following
// interface:
//
//     type HasUnderlying interface {
//            Cause() error
//     }
//
// If the error does not implement Cause, the original error will
// be returned. If the error is nil, nil will be returned without further
// investigation.
func Cause(err error) error {
	for err != nil {
		cause, ok := err.(HasUnderlying)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return err
}

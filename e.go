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

// E is error with more information. It is able to marshall itself to json as response.
// Result of Error() method should include stacktrace. Therefore it should not be
// displayed directly to the user
type E interface {
	error
	json.Marshaler
	IsReq() bool
	WithMsg(string) E
	Stacktrace() *stack.Multi
}

type errstack struct {
	err        error
	stacktrace *stack.Multi
	message    string
}

type err errstack

func newErr(e error, s string, skip int) err {
	st := stack.CallersMulti(skip + 1)
	return err{e, st, s}
}

func (e *err) Error() string {
	if e.message == "" && e.err == nil {
		return "error"
	}
	if e.err == nil {
		return e.message
	}
	if e.message == "" {
		return e.err.Error()
	}
	return fmt.Sprintf("%s [%s]", e.message, e.err.Error())
}

func (e err) withMsg(msg string) err {
	return err{
		err:        wrapper{e.message, e.err},
		message:    msg,
		stacktrace: e.stacktrace,
	}
}

func (e *err) Cause() error {
	return Cause(e.err)
}

func (e *err) Stacktrace() *stack.Multi {
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

package errstack

import (
	"encoding/json"
	"fmt"

	"github.com/facebookgo/stack"
)

// E is error with more information. It is able to marshall itself to json as response.
// Result of Error() method should include stacktrace. Therefore it should not be
// displayed directly to the user
type E interface {
	error
	json.Marshaler
	Inf() bool
	WithMsg(string) E
	Stacktrace() *stack.Multi
}

// HasUnderlying describes entity (usually an error) which has underlying error.
type HasUnderlying interface {
	Underlying() error
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

func newErrF(e error, skip int, f string, a ...interface{}) err {
	return newErr(e, fmt.Sprintf(f, a...), skip+1)
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

func (e *err) Underlying() error {
	return underlying(e.err)
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

func (e *wrapper) Underlying() error {
	return underlying(e.err)
}

func underlying(err error) error {
	if hasUnderlying, ok := err.(HasUnderlying); ok {
		return hasUnderlying.Underlying()
	}
	return err
}

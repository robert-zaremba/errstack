package errstack

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/facebookgo/stack"
)

type errstack struct {
	err        error
	stacktrace stack.Stack
	msg        string
	kind       Kind
	details    map[string]interface{}
}

func newErr(e error, s string, kind Kind, skip int) *errstack {
	st := stack.Callers(skip + 1)
	return &errstack{e, st, s, kind, map[string]interface{}{}}
}

// New creates a new error E
func New(kind Kind, msg string) E {
	return newErr(nil, msg, kind, 1)
}

// Wrap creates new error using error and string message
func Wrap(err error, kind Kind, msg string) E {
	if err == nil {
		return nil
	}
	if es, ok := err.(errstack); ok && es.kind == kind {
		return es.WithMsg(msg)
	}
	return newErr(err, msg, kind, 1)
}

func (e errstack) WithMsg(msg string) E {
	return errstack{
		err:        wrapper{e.msg, e.err},
		msg:        msg,
		stacktrace: e.stacktrace,
		kind:       e.kind,
	}
}

// IsReq is false for Infrastructure errors.
// It implements errstack.E interface
func (e errstack) IsReq() bool {
	return isReq(e.kind)
}

// Kind implements E interface.
// returns error kind type.
func (e errstack) Kind() Kind {
	return e.kind
}

// Details implements E interface.
func (e errstack) Details() map[string]interface{} {
	return e.details
}

// Add implements E interface.
func (e errstack) Add(key string, payload interface{}) {
	e.details[key] = payload
}

// StatusCode return HTTP status code
func (e errstack) StatusCode() int {
	if s, ok := e.err.(HasStatusCode); ok {
		return s.StatusCode()
	}
	if e.IsReq() {
		return 400
	}
	return 500
}

func (e errstack) Cause() error {
	return e.err
}

// Stacktrace returns error creation stacktrace
func (e errstack) Stacktrace() stack.Stack {
	return e.stacktrace
}

func (e errstack) Error() string {
	if e.err == nil {
		return e.msg
	}
	return fmt.Sprint(e.msg, " [", e.err.Error(), "]")
}

// MarshalJSON implements Marshaller
// It will return "Internal server error" without full details when the error
// is not a request error.
func (e errstack) MarshalJSON() ([]byte, error) {
	if e.IsReq() {
		data := errmap{"msg": e.msg}
		if e.err != nil {
			if _, ok := e.err.(json.Marshaler); ok {
				data["err"] = e.err
			} else {
				data["err"] = e.err.Error()
			}
		}
		return json.Marshal(data)
	}
	return json.Marshal("Internal server error: " + e.msg)
}

// Format implements fmt.Formatter interface
func (e errstack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			io.WriteString(s, "### ")
			io.WriteString(s, e.Error())
			io.WriteString(s, "\n")
			io.WriteString(s, e.stacktrace.String())
			io.WriteString(s, "\n--------------------------------")
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, e.msg)
	case 'q':
		fmt.Fprintf(s, "%q", e.msg)
	}
}

type wrapper struct {
	msg string
	err error
}

func (e wrapper) Error() string {
	if e.err == nil {
		return e.msg
	}
	return fmt.Sprintf("%s [%s]", e.msg, e.err.Error())
}

func (e wrapper) MarshalJSON() ([]byte, error) {
	data := errmap{"msg": e.msg}
	if e.err != nil {
		data["err"] = e.err
	}
	return json.Marshal(data)
}

func (e wrapper) Cause() error {
	return e.err
}

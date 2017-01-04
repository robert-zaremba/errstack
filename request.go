package errstack

import (
	"encoding/json"

	"github.com/davecgh/go-spew/spew"
	"github.com/facebookgo/stack"
)

type request struct {
	details    errmap
	stacktrace *stack.Multi
}

func (r *request) msg() string {
	if len(r.details) == 0 {
		return "requst error"
	}
	return spew.Sdump(r.details)
}

// Error implements error interface
func (r *request) Error() string {
	return r.msg()
}

// Inf implements errstack.E interface
func (r *request) Inf() bool {
	return false
}

func (r *request) Stacktrace() *stack.Multi {
	return r.stacktrace
}

// MarshalJSON implements Marshaller interface
func (r *request) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.details)
}

func (r *request) WithMsg(msg string) E {
	return &simpleRequest{err{
		err:        r.details,
		message:    msg,
		stacktrace: r.stacktrace,
	}}
}

func newRequest(m map[string]interface{}, skip int) E {
	st := stack.CallersMulti(skip + 1)
	return &request{m, st}
}

// NewReqDetails creates a request error.
// Key inform which request parameter was invalid
// and details contains reason of error
func NewReqDetails(key string, details interface{}) E {
	return newRequest(errmap{key: details}, 1)
}

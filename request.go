package errstack

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/davecgh/go-spew/spew"
	"github.com/facebookgo/stack"
)

type request struct {
	details    errmap
	msg        string
	stacktrace stack.Stack
}

func init() {
	// Fixes error output when printing request.details
	spew.Config.DisableMethods = true
}

// IsReq is true for Request errors.
// It implements errstack.E interface
func (r *request) IsReq() bool {
	return true
}

// Kind implements E interface.
// returns error kind type.
func (r *request) Kind() Kind {
	return Request
}

// Details implements E interface.
func (r *request) Details() map[string]interface{} {
	return r.details
}

// Add implements E interface.
func (r *request) Add(key string, payload interface{}) {
	r.details[key] = payload
}

// StatusCode return HTTP status code
func (r *request) StatusCode() int {
	return 400
}

func (r *request) Stacktrace() stack.Stack {
	return r.stacktrace
}

// Error implements error interface
func (r *request) Error() string {
	return fmt.Sprintf("%s %v", r.msg, r.details)
}

// MarshalJSON implements Marshaller interface
func (r *request) MarshalJSON() ([]byte, error) {
	if r.msg == "" {
		return json.Marshal(r.details)
	}
	return json.Marshal(errmap{"msg": r.msg, "err": r.details})
}

// Format implements fmt.Formatter interface
func (r *request) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			io.WriteString(s, "### ")
			io.WriteString(s, r.msg)
			fmt.Fprintf(s, " %v\n", r.details)
			io.WriteString(s, r.stacktrace.String())
			io.WriteString(s, "\n--------------------------------")
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, r.msg)
	case 'q':
		fmt.Fprintf(s, "%q", r.msg)
	}
}

func (r *request) WithMsg(msg string) E {
	r2 := *r // make a copy
	r2.msg = fmt.Sprintf("%s [%s]", msg, r.msg)
	return &r2
}

func newRequest(m map[string]interface{}, msg string, skip int) E {
	st := stack.Callers(skip + 1)
	return &request{m, msg, st}
}

// NewReqDetails creates a request error.
// Key inform which request parameter was invalid
// and details contains reason of error
func NewReqDetails(key string, details interface{}, msg string) E {
	return newRequest(errmap{key: details}, msg, 1)
}

package errstack

import (
	"encoding/json"
	"fmt"
)

type simpleRequest struct {
	err
}

func (r *simpleRequest) WithMsg(msg string) E {
	return &simpleRequest{r.withMsg(msg)}
}

// Inf implements errstack.E interface
func (r *simpleRequest) Inf() bool {
	return false
}

func newSimpleRequest(e error, d string, skip int) E {
	return &simpleRequest{newErr(e, d, skip+1)}
}

func (r *simpleRequest) MarshalJSON() ([]byte, error) {
	data := errmap{"msg": r.Error()}
	return json.Marshal(data)
}

// NewReqF creates request error from format
func NewReqF(f string, a ...interface{}) E {
	return newSimpleRequest(nil, fmt.Sprintf(f, a...), 1)
}

// NewReq creates request error from string
func NewReq(s string) E {
	return newSimpleRequest(nil, s, 1)
}

// WrapAsReq creates new request error from simple error
// If input argument is nil, nil is returned.
// If input argument is already errstack.E, it is returned unchanged.
func WrapAsReq(err error, message string) E {
	return wrapSimpleRequest(err, message, 1)
}

// WrapAsReqF creates new request error from simple error and creates message from format
// If input argument is nil, nil is returned.
// If input argument is already errstack.E, it is returned unchanged.
func WrapAsReqF(err error, f string, a ...interface{}) E {
	return wrapSimpleRequest(err, fmt.Sprintf(f, a...), 1)
}

func wrapSimpleRequest(err error, message string, skip int) E {
	if err == nil {
		return nil
	}
	switch ee := err.(type) {
	case *simpleRequest:
		return ee.WithMsg(message)
	case *request:
		return ee.WithMsg(message)
	case E:
		// TODO: solve the reporting way.
		// logger.Warn("Can't wrap non-request errstack.E as request", "stack", stack.Callers(skip+1))
		return ee
	}
	return newSimpleRequest(err, message, skip+1)
}

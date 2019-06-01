package errstack

import (
	"fmt"
)

func wrapSimpleRequest(e error, message string, skip int) E {
	if e == nil {
		return nil
	}
	switch ee := e.(type) {
	case *request:
		return ee.WithMsg(message)
	case errstack:
		if ee.kind == Request {
			return ee.WithMsg(message)
		}
	}
	return newErr(e, message, Request, skip+1)
}

// NewReqF creates request error from format
func NewReqF(f string, a ...interface{}) E {
	return newErr(nil, fmt.Sprintf(f, a...), Request, 1)
}

// NewReq creates request error from string
func NewReq(s string) E {
	return newErr(nil, s, Request, 1)
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

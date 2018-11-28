package errstack

import (
	"encoding/json"
	"fmt"
)

type infrastructure struct {
	errstack
}

func (i *infrastructure) WithMsg(msg string) E {
	return &infrastructure{i.withMsg(msg)}
}

// IsReq is false for Infrastructure errors.
// It implements errstack.E interface
func (i *infrastructure) IsReq() bool {
	return false
}

// StatusCode return HTTP status code
func (i *infrastructure) StatusCode() int {
	if s, ok := i.errstack.err.(HasStatusCode); ok {
		return s.StatusCode()
	}
	return 500
}

// MarshalJSON implements Marshaller
func (i *infrastructure) MarshalJSON() ([]byte, error) {
	if m, ok := i.errstack.err.(json.Marshaler); ok {
		return m.MarshalJSON()
	}
	return []byte(`"Internal server error"`), nil
}

func newInfrastructure(details string, skip int) E {
	return &infrastructure{newErr(nil, details, skip+1)}
}

func wrapInfrastructure(e error, details string, skip int) E {
	if e == nil {
		return nil
	}
	if errstack, ok := e.(E); ok {
		return errstack
	}
	return &infrastructure{newErr(e, details, skip+1)}
}

// WrapAsInf creates new infrastructure error from simple error
// If input argument is nil, nil is returned.
// If input argument is already errstack.E, it is returned unchanged.
func WrapAsInf(e error, messages ...string) E {
	var msg string
	if len(messages) != 0 {
		msg = messages[0]
	}
	return wrapInfrastructure(e, msg, 1)
}

// WrapAsInfF creates new infrastructural error wrapping given error and
// using string formatter for description.
func WrapAsInfF(err error, f string, a ...interface{}) E {
	return wrapInfrastructure(err, fmt.Sprintf(f, a...), 1)
}

// NewInfF creates new infrastructural error using string formatter
func NewInfF(f string, a ...interface{}) E {
	return newInfrastructure(fmt.Sprintf(f, a...), 1)
}

// NewInf creates new infrastructural error from string
func NewInf(s string) E {
	return newInfrastructure(s, 1)
}

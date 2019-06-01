package errstack

import (
	"fmt"
)

func wrapIO(e error, details string, skip int) E {
	if e == nil {
		return nil
	}
	if es, ok := e.(errstack); ok && es.kind == IO {
		return es.WithMsg(details)
	}
	return newErr(e, details, IO, skip+1)
}

// WrapAsIO creates new infrastructure error from simple error
// If input argument is nil, nil is returned.
func WrapAsIO(e error, messages ...string) E {
	var msg string
	if len(messages) != 0 {
		msg = messages[0]
	}
	return wrapIO(e, msg, 1)
}

// WrapAsIOf creates new infrastructural error wrapping given error and
// using string formatter for description.
func WrapAsIOf(err error, f string, a ...interface{}) E {
	return wrapIO(err, fmt.Sprintf(f, a...), 1)
}

// NewIOf creates new infrastructural error using string formatter
func NewIOf(f string, a ...interface{}) E {
	return newErr(nil, fmt.Sprintf(f, a...), IO, 1)
}

// NewIO creates new infrastructural error from string
func NewIO(s string) E {
	return newErr(nil, s, IO, 1)
}

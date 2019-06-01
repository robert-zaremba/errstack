package errstack

import "fmt"

func newDomain(details string, skip int) E {
	return newErr(nil, details, Domain, skip+1)
}

func wrapDomain(e error, details string, skip int) E {
	if e == nil {
		return nil
	}
	if es, ok := e.(errstack); ok && es.kind == Domain {
		return es.WithMsg(details)
	}
	return newErr(e, details, Domain, skip+1)
}

// NewDomainF creates new domain error using string formatter
func NewDomainF(format string, a ...interface{}) E {
	return newDomain(fmt.Sprintf(format, a...), 1)
}

// NewDomain creates new domain error from string
// Domain error is classified as an Infrastructure error.
func NewDomain(s string) E {
	return newDomain(s, 1)
}

// WrapAsDomain creates new domain error using error and string message
// Domain error is classified as an Infrastructure error.
func WrapAsDomain(err error, message string) E {
	return wrapDomain(err, message, 1)
}

// WrapAsDomainF creates new domain error wrapping given error and
// using string formatter for description.
// Domain error is classified as an Infrastructure error.
func WrapAsDomainF(err error, f string, a ...interface{}) E {
	return wrapDomain(err, fmt.Sprintf(f, a...), 1)
}

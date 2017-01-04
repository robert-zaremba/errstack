package errstack

import "fmt"

// NewDomainF creates new domain error using string formatter
// domain error represent an error when our model is in wrong state
// (eg user is in impossible state from our domain stand). Compared to infrastructure
// error - when it is related that some operation failed on other service.
func NewDomainF(format string, a ...interface{}) E {
	return newInfrastructure(fmt.Sprintf(format, a...), 1)
}

// NewDomain creates new domain error from string
func NewDomain(s string) E {
	return newInfrastructure(s, 1)
}

// WrapAsDomain creates new domain error using error and string message
func WrapAsDomain(err error, message string) E {
	return wrapInfrastructure(err, message, 1)
}

// WrapAsDomainF creates new domain error wrapping given error and
// using string formatter for description.
func WrapAsDomainF(err error, f string, a ...interface{}) E {
	return wrapInfrastructure(err, fmt.Sprintf(f, a...), 1)
}

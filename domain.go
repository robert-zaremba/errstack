package errstack

import "fmt"

// Domain error represent an error when our model is in wrong state
// (eg user is in impossible state from our domain stand). Compared to infrastructure
// error - when it is related that some operation failed on other service.
// Domain error is classified as a non Request error.
type Domain *infrastructure

// NewDomainF creates new domain error using string formatter
func NewDomainF(format string, a ...interface{}) E {
	return newInfrastructure(fmt.Sprintf(format, a...), 1)
}

// NewDomain creates new domain error from string
// Domain error is classified as an Infrastructure error.
func NewDomain(s string) E {
	return newInfrastructure(s, 1)
}

// WrapAsDomain creates new domain error using error and string message
// Domain error is classified as an Infrastructure error.
func WrapAsDomain(err error, message string) E {
	return wrapInfrastructure(err, message, 1)
}

// WrapAsDomainF creates new domain error wrapping given error and
// using string formatter for description.
// Domain error is classified as an Infrastructure error.
func WrapAsDomainF(err error, f string, a ...interface{}) E {
	return wrapInfrastructure(err, fmt.Sprintf(f, a...), 1)
}

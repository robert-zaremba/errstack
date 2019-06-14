package errstack

import (
	"encoding/json"

	"github.com/facebookgo/stack"
)

// HasUnderlying describes entity (usually an error) which has underlying error.
type HasUnderlying interface {
	// Cause returns the underlying error (if any) causing the failure.
	Cause() error
}

// HasStatusCode provides a function to return the HTTP status code.
type HasStatusCode interface {
	StatusCode() int
}

// HasStacktrace provides a function to return the the root stacktrace
type HasStacktrace interface {
	Stacktrace() stack.Stack
}

// E is error with more information. It is able to marshall itself to json as response.
// Result of Error() method should include stacktrace. Therefore it should not be
// displayed directly to the user
type E interface {
	error
	HasStatusCode
	HasStacktrace
	json.Marshaler
	IsReq() bool
	Kind() Kind
	WithMsg(string) E
	Details() map[string]interface{}
	Add(key string, payload interface{}) // add more details to the error
}

// Kind defines the kind of error that must act differently depending on the error
type Kind uint8

// Kinds of errors.
const (
	Other         Kind = iota // Unclassified error.
	Invalid                   // Invalid operation for this type of item.
	Permission                // Permission denied.
	IO                        // I/O error such as network failure.
	Exist                     // Item already exists.
	NotExist                  // Item does not exist.
	IsDir                     // Item is a directory.
	NotDir                    // Item is not a directory.
	NotEmpty                  // Directory not empty.
	Private                   // Information withheld.
	CannotDecrypt             // No wrapped key for user with read access.
	Transient                 // A transient error.
	BrokenLink                // Link target does not exist.
	Request                   // General request error
	Domain                    // Internal error causing business domain problem or inconsistency
)

// IsKind reports whether err is an *Error of the given Kind.
// If err is nil then Is returns false.
func IsKind(kind Kind, err error) bool {
	e, ok := err.(E)
	if !ok {
		return false
	}
	ek := e.Kind()
	if ek != Other {
		return ek == kind
	}
	if ecauser, ok := err.(HasUnderlying); ok {
		err = ecauser.Cause()
		if err != nil {
			return IsKind(kind, err)
		}
	}
	return false
}

func isReq(kind Kind) bool {
	return kind == Permission || kind == Exist || kind == NotExist || kind == Private || kind == CannotDecrypt || kind == Request
}

// RootErr returns the underlying cause of the error, if possible.
// Normally it should be the root error.
// This method uses `HasUnderlying` interface to extract the cause error.
//
// If the error does not implement the `HasUnderlying`, the original error will
// be returned. If the error is nil, nil will be returned without further
// investigation.
func RootErr(err error) error {
	for err != nil {
		cause, ok := err.(HasUnderlying)
		if !ok {
			break
		}
		errC := cause.Cause()
		if errC == nil {
			return err
		}
		err = errC
	}
	return err
}

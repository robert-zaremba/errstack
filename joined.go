package errstack

import "fmt"

type joinedError struct {
	errors []error
}

func (je joinedError) Error() string {
	return fmt.Sprint("<JoinedError ", je.errors, ">")
}

// Join creates a new error from list of errors. It filters out nil errors.
// If there is no not-nil error it returs nil.
func Join(es ...error) error {
	var filtered = []error{}
	for _, e := range es {
		if e != nil {
			filtered = append(filtered, e)
		}
	}
	if len(filtered) == 0 {
		return nil
	}
	return joinedError{filtered}
}

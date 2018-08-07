package errors

import "github.com/pkg/errors"

// ExitError is an error with an embedded exit status.
type ExitError struct {
	Err      error
	ExitCode int
}

func (e *ExitError) Error() string {
	return e.Err.Error()
}

// ExitStatus returns the error's embedded status
func (e *ExitError) ExitStatus() int {
	return e.ExitCode
}

// Cause returns the original cause.
func (e *ExitError) Cause() error {
	return e.Err
}

// Errorf calls pkg/error's Errorf function.
func Errorf(format string, args ...interface{}) error {
	return errors.Errorf(format, args...)
}

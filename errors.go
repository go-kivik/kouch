package kouch

import (
	"fmt"
	"os"
)

type exitStatuser interface {
	ExitStatus() int
}

// ExitStatus returns the exit status embedded in the error.
func ExitStatus(err error) int {
	if err == nil {
		return 0
	}
	if statuser, ok := err.(exitStatuser); ok {
		return statuser.ExitStatus()
	}
	return ExitUnknownFailure
}

// Exit outputs err.Error() to stderr, then exits with the exit status embedded
// in the error.
func Exit(err error) {
	fmt.Fprintf(os.Stderr, err.Error())
	os.Exit(ExitStatus(err))
}

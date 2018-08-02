package kouch

import (
	"fmt"
	"os"

	"github.com/go-kivik/kivik"
)

type exitStatuser interface {
	ExitStatus() int
}

// ExitStatus returns the exit status embedded in the error.
func ExitStatus(err error) int {
	if err == nil {
		return 0
	}
	if statuser, ok := err.(exitStatuser); ok { // nolint: misspell
		return statuser.ExitStatus()
	}
	return ExitUnknownFailure
}

// Exit outputs err.Error() to stderr, then exits with the exit status embedded
// in the error.
func Exit(err error) {
	exitStatus := ExitStatus(err)
	httpStatus := kivik.StatusCode(err)
	if httpStatus >= 400 && httpStatus < 600 {
		fmt.Fprintf(os.Stderr, "kouch: (%d) The requested URL returned error: %d %s\n",
			exitStatus, httpStatus, err)
	} else {
		fmt.Fprintf(os.Stderr, "kouch: (%d) %s\n", exitStatus, err)
	}
	os.Exit(exitStatus)
}

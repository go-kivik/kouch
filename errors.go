package kouch

import (
	"fmt"
	"os"

	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kivik"
)

// InitError returns an error for init failures.
type InitError string

func (i InitError) Error() string { return string(i) }

// ExitStatus returns ExitFailedToInitialize
func (i InitError) ExitStatus() int { return chttp.ExitFailedToInitialize }

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
	return chttp.ExitUnknownFailure
}

// Exit outputs err.Error() to stderr, then exits with the exit status embedded
// in the error.
func Exit(err error) {
	msg, exitStatus := exit(err)
	_, _ = fmt.Fprintf(os.Stderr, "kouch: (%d) %s\n", exitStatus, msg)
	os.Exit(exitStatus)
}

func exit(err error) (string, int) {
	exitStatus := ExitStatus(err)
	httpStatus := kivik.StatusCode(err)
	if exitStatus == chttp.ExitNotRetrieved && httpStatus >= 400 && httpStatus < 600 {
		return fmt.Sprintf("The requested URL returned error: %d %s", httpStatus, err), exitStatus
	}
	return err.Error(), exitStatus
}

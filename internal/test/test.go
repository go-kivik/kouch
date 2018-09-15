// Package test provides test utilities
package test

import (
	"os"
	"testing"

	"github.com/flimzy/diff"
	"github.com/flimzy/testy"
	"github.com/go-kivik/kouch/cmd/kouch/registry"
)

// CmdTest represents a single test for a complete kouch command.
type CmdTest struct {
	Args   []string
	Stdout string
	Stderr string
	Err    string
	Status int
}

// ValidateCmdTest validates a test. It is meant to be passed to
// testy's tests.Run() method.
func ValidateCmdTest(args []string) func(*testing.T, CmdTest) {
	return func(t *testing.T, test CmdTest) {
		defer testy.RestoreEnv()()
		if e := os.Setenv("HOME", "/dev/null"); e != nil {
			t.Fatal(e)
		}
		var err error
		stdout, stderr := testy.RedirIO(nil, func() {
			root := registry.Root()
			root.SetArgs(append(args, test.Args...))
			err = root.Execute()
		})
		if d := diff.Text(test.Stdout, stdout); d != nil {
			t.Errorf("STDOUT:\n%s", d)
		}
		if d := diff.Text(test.Stderr, stderr); d != nil {
			t.Errorf("STDERR:\n%s", d)
		}
		testy.ExitStatusError(t, test.Err, test.Status, err)
	}
}

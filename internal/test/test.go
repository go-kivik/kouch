// Package test provides test utilities
package test

import (
	"testing"

	"github.com/flimzy/diff"
	"github.com/flimzy/testy"
	"github.com/go-kivik/kouch"
	"github.com/go-kivik/kouch/cmd/kouch/registry"
)

// CmdTest represents a single test for a complete kouch command.
type CmdTest struct {
	Conf   *kouch.Config
	Args   []string
	Stdout string
	Stderr string
	Err    string
	Status int
}

// ValidateCmdTest validates a test. It is meant to be passed to
// testy's tests.Run() method.
func ValidateCmdTest(t *testing.T, test CmdTest) {
	var err error
	stdout, stderr := testy.RedirIO(nil, func() {
		root := registry.Root()
		root.SetArgs(append([]string{"get", "uuids"}, test.Args...))
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

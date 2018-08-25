package testy

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/flimzy/diff"
)

func TestRedirIO(t *testing.T) {
	tests := []struct {
		name           string
		fn             func()
		stdin          io.Reader
		stdout, stderr string
	}{
		{
			name:   "stdout",
			fn:     func() { fmt.Printf("some output\n") },
			stdout: "some output\n",
		},
		{
			name:   "stderr",
			fn:     func() { fmt.Fprintf(os.Stderr, "some error output\n") },
			stderr: "some error output\n",
		},
		{
			name: "echo stdin to stdout",
			fn: func() {
				io.Copy(os.Stdout, os.Stdin)
			},
			stdin:  strings.NewReader("testing"),
			stdout: "testing",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			out, err := RedirIO(test.stdin, test.fn)
			if d := diff.Text(test.stdout, out); d != nil {
				t.Errorf("stdout:\n%s", d)
			}
			if d := diff.Text(test.stderr, err); d != nil {
				t.Errorf("stderr:\n%s", d)
			}
			if os.Stdin.Fd() != 0 {
				t.Errorf("STDIN not restored properly\n")
			}
			if os.Stdout.Fd() != 1 {
				t.Errorf("STDOUT not restored properly\n")
			}
			if os.Stderr.Fd() != 2 {
				t.Errorf("STDERR not restored properly\n")
			}
		})
	}
}

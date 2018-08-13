package root

import (
	"testing"

	"github.com/flimzy/testy"
	"github.com/go-kivik/kouch"
)

func TestVerbose(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected bool
		err      string
	}{
		{
			name:     "defaults",
			expected: false,
		},
		{
			name:     "verbose enabled",
			args:     []string{"--" + flagVerbose},
			expected: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cmd := rootCmd("1.2.3")
			cmd.ParseFlags(test.args)
			ctx, err := verbose(kouch.GetContext(cmd), cmd)
			testy.Error(t, test.err, err)
			if verbose := kouch.Verbose(ctx); verbose != test.expected {
				t.Errorf("Unexpected result: %t\n", verbose)
			}
		})
	}
}

package kouch

import (
	"os"
	"testing"

	"github.com/flimzy/testy"
)

func TestHome(t *testing.T) {
	tests := []struct {
		name     string
		env      map[string]string
		expected string
	}{
		{
			name:     "No HOME",
			expected: "",
		},
		{
			name:     "Configured home",
			env:      map[string]string{"HOME": "/home/joe"},
			expected: "/home/joe/.kouch",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer testy.RestoreEnv()()
			os.Clearenv()
			testy.SetEnv(test.env)
			home := Home()
			if home != test.expected {
				t.Errorf("Unexpected home result: '%s', expected '%s'\n", home, test.expected)
			}
		})
	}
}

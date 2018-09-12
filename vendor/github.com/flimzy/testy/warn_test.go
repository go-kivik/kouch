package testy

import (
	"regexp"
	"testing"
)

func TestSwarn(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		args     []interface{}
		expected string
	}{
		{
			name:     "simple",
			format:   "foo %s",
			args:     []interface{}{"bar"},
			expected: `^\[\d+/\d+\ .*/github.com/flimzy/testy/warn_test.go:24] foo bar$`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Swarn(test.format, test.args...)
			if !regexp.MustCompile(test.expected).MatchString(result) {
				t.Errorf("Unexpected result: %s\nExpected: /%s/\n", result, test.expected)
			}
		})
	}
}

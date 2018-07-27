package io

import (
	"bytes"
	"testing"

	"github.com/flimzy/diff"
	"github.com/flimzy/testy"
)

func TestYAMLOutput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		err      string
	}{
		{
			name:  "happy path",
			input: `{"foo":"bar", "baz":123, "qux": [1,2,3]}`,
			expected: `baz: 123
foo: bar
qux:
- 1
- 2
- 3`,
		},
		{
			name:  "invalid JSON input",
			input: "oink",
			err:   `invalid character 'o' looking for beginning of value`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := &yamlProcessor{}
			buf := &bytes.Buffer{}
			err := p.Output(buf, []byte(test.input))
			testy.Error(t, test.err, err)
			if d := diff.Text(test.expected, buf.String()); d != nil {
				t.Error(d)
			}
		})
	}
}

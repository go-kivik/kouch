package io

import (
	"bytes"
	"testing"

	"github.com/flimzy/diff"
	"github.com/flimzy/testy"
)

func TestJSONOutput(t *testing.T) {
	tests := []struct {
		name           string
		prefix, indent string
		escapeHTML     bool
		input          string
		expected       string
		err            string
	}{
		{
			name:     "happy path",
			input:    `{"foo":"bar", "baz":123}`,
			expected: `{"baz":123,"foo":"bar"}`,
		},
		{
			name:   "indented",
			indent: "\t",
			input:  `{"foo":"bar", "baz":123}`,
			expected: `{
	"baz": 123,
	"foo": "bar"
}`,
		},
		{
			name:       "no escape HTML",
			escapeHTML: false,
			input:      `{"foo": "<>"}`,
			expected:   `{"foo":"<>"}`,
		},
		{
			name:       "escape HTML",
			escapeHTML: true,
			input:      `{"foo": "<>"}`,
			expected:   `{"foo":"\u003c\u003e"}`,
		},
		{
			name:  "invalid JSON input",
			input: "oink",
			err:   `invalid character 'o' looking for beginning of value`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := &jsonProcessor{
				prefix:     test.prefix,
				indent:     test.indent,
				escapeHTML: test.escapeHTML,
			}
			buf := &bytes.Buffer{}
			err := p.Output(buf, []byte(test.input))
			testy.Error(t, test.err, err)
			if d := diff.Text(test.expected, buf.String()); d != nil {
				t.Error(d)
			}
		})
	}
}

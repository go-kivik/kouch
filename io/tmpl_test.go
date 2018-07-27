package io

import (
	"bytes"
	"testing"

	"github.com/flimzy/diff"
	"github.com/flimzy/testy"
)

func TestTmplOutput(t *testing.T) {
	tests := []struct {
		name     string
		template string
		input    string
		expected string
		err      string
	}{
		{
			name:     "happy path",
			template: `{{ .foo }}`,
			input:    `{"foo":"bar", "baz":123, "qux": [1,2,3]}`,
			expected: `bar`,
		},
		{
			name:  "invalid JSON input",
			input: "oink",
			err:   `invalid character 'o' looking for beginning of value`,
		},
		{
			name:     "invalid template",
			template: "{{ .foo }",
			input:    `{"foo":"bar"}`,
			err:      `template: :1: unexpected "}" in operand`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := &tmplProcessor{template: test.template}
			buf := &bytes.Buffer{}
			err := p.Output(buf, []byte(test.input))
			testy.Error(t, test.err, err)
			if d := diff.Text(test.expected, buf.String()); d != nil {
				t.Error(d)
			}
		})
	}
}

package io

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/flimzy/diff"
	"github.com/flimzy/testy"
	"github.com/go-kivik/kouch/kouchio"
	"github.com/spf13/cobra"
)

func TestJsonModeConfig(t *testing.T) {
	cmd := &cobra.Command{}
	mode := &jsonMode{}
	mode.AddFlags(cmd.PersistentFlags())

	testOptions(t, []string{"json-escape-html", "json-indent", "json-prefix"}, cmd)
}

func TestJSONOutput(t *testing.T) {
	tests := []struct {
		name             string
		args             []string
		flagsErr, newErr string
		input            string
		expected         string
		err              string
	}{
		{
			name:     "happy path",
			input:    `{"foo":"bar", "baz":123}`,
			expected: `{"baz":123,"foo":"bar"}`,
		},
		{
			name:  "indented",
			args:  []string{"--json-indent", "\t"},
			input: `{"foo":"bar", "baz":123}`,
			expected: `{
	"baz": 123,
	"foo": "bar"
}`,
		},
		{
			name:  "prefix",
			args:  []string{"--json-prefix", "xx"},
			input: `{"foo":"bar", "baz":123}`,
			expected: `{
xx"baz": 123,
xx"foo": "bar"
xx}`,
		},
		{
			name:     "no escape HTML",
			input:    `{"foo": "<>"}`,
			expected: `{"foo":"<>"}`,
		},
		{
			name:     "escape HTML",
			args:     []string{"--json-escape-html"},
			input:    `{"foo": "<>"}`,
			expected: `{"foo":"\u003c\u003e"}`,
		},
		{
			name:  "invalid JSON input",
			input: "oink",
			err:   `invalid character 'o' looking for beginning of value`,
		},
		{
			name:     "invalid args",
			args:     []string{"--foo"},
			flagsErr: "unknown flag: --foo",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			mode := &jsonMode{}
			mode.AddFlags(cmd.PersistentFlags())

			err := cmd.ParseFlags(test.args)
			testy.Error(t, test.flagsErr, err)

			buf := &bytes.Buffer{}
			p, err := mode.New(cmd.Flags(), buf)
			testy.Error(t, test.newErr, err)

			defer kouchio.CloseWriter(p) // nolint: errcheck
			_, err = io.Copy(p, strings.NewReader(test.input))
			if err == nil {
				err = kouchio.CloseWriter(p)
			}
			testy.Error(t, test.err, err)
			if d := diff.Text(test.expected, buf.String()); d != nil {
				t.Error(d)
			}
		})
	}
}

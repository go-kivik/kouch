package io

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/flimzy/diff"
	"github.com/flimzy/testy"
	"github.com/spf13/cobra"
)

func TestJsonModeConfig(t *testing.T) {
	cmd := &cobra.Command{}
	mode := &jsonMode{}
	mode.config(cmd.PersistentFlags())

	testOptions(t, []string{"json-escape-html", "json-indent", "json-prefix"}, cmd)
}

func TestJsonNew(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		parseErr string
		expected *jsonProcessor
		err      string
	}{
		{
			name:     "happy path, no options",
			args:     nil,
			expected: &jsonProcessor{},
		},
		{
			name:     "invalid args",
			args:     []string{"--foo"},
			parseErr: "unknown flag: --foo",
		},
		{
			name: "happy path, prefix",
			args: []string{"--json-prefix", "xx"},
			expected: &jsonProcessor{
				prefix: "xx",
			},
		},
		{
			name: "happy path, indent",
			args: []string{"--json-indent", "--"},
			expected: &jsonProcessor{
				indent: "--",
			},
		},
		{
			name: "happy path, escape html",
			args: []string{"--json-escape-html"},
			expected: &jsonProcessor{
				escapeHTML: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			mode := &jsonMode{}
			mode.config(cmd.PersistentFlags())

			err := cmd.ParseFlags(test.args)
			testy.Error(t, test.parseErr, err)

			result, err := mode.new(cmd)
			testy.Error(t, test.err, err)
			if d := diff.Interface(test.expected, result); d != nil {
				t.Error(d)
			}
		})
	}
}

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
			err := p.Output(buf, ioutil.NopCloser(strings.NewReader(test.input)))
			testy.Error(t, test.err, err)
			if d := diff.Text(test.expected, buf.String()); d != nil {
				t.Error(d)
			}
		})
	}
}

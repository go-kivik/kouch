package io

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/flimzy/diff"
	"github.com/flimzy/testy"
	"github.com/spf13/cobra"
)

func TestYamlModeConfig(t *testing.T) {
	cmd := &cobra.Command{}
	mode := &yamlMode{}
	mode.config(cmd.PersistentFlags())

	testOptions(t, []string{}, cmd)
}

func TestYamlNew(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		parseErr string
		expected *yamlProcessor
		err      string
	}{
		{
			name:     "happy path",
			args:     nil,
			expected: &yamlProcessor{underlying: &bytes.Buffer{}},
		},
		{
			name:     "invalid args",
			args:     []string{"--foo"},
			parseErr: "unknown flag: --foo",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			mode := &yamlMode{}
			mode.config(cmd.PersistentFlags())

			err := cmd.ParseFlags(test.args)
			testy.Error(t, test.parseErr, err)

			buf := &bytes.Buffer{}
			result, err := mode.new(cmd, buf)
			testy.Error(t, test.err, err)
			if d := diff.Interface(test.expected, result); d != nil {
				t.Error(d)
			}
		})
	}
}

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
			buf := &bytes.Buffer{}
			p := newYAMLProcessor(buf)
			defer p.Close()
			_, err := io.Copy(p, strings.NewReader(test.input))
			if err == nil {
				err = p.Close()
			}
			testy.Error(t, test.err, err)
			if d := diff.Text(test.expected, buf.String()); d != nil {
				t.Error(d)
			}
		})
	}
}

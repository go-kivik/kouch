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

func TestYAMLOutput(t *testing.T) {
	tests := []struct {
		name             string
		cmd              *cobra.Command
		args             []string
		flagsErr, newErr string
		input            string
		expected         string
		err              string
	}{
		{
			name:     "invalid args",
			args:     []string{"--foo"},
			flagsErr: "unknown flag: --foo",
		},
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
			cmd := &cobra.Command{}
			mode := &yamlMode{}
			mode.config(cmd.PersistentFlags())

			err := cmd.ParseFlags(test.args)
			testy.Error(t, test.flagsErr, err)

			buf := &bytes.Buffer{}
			p, err := mode.new(cmd.Flags(), buf)
			testy.Error(t, test.newErr, err)

			defer p.Close() // nolint: errcheck
			_, err = io.Copy(p, strings.NewReader(test.input))
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

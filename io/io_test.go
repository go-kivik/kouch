package io

import (
	"testing"

	"github.com/flimzy/diff"
	"github.com/flimzy/testy"
	"github.com/spf13/cobra"
)

func TestAddFlags(t *testing.T) {
	cmd := &cobra.Command{}
	AddFlags(cmd)

	testOptions(t, []string{"json-escape-html", "json-indent", "json-prefix", "output-format", "template", "template-file"}, cmd)
}

func TestSelectOutputProcessor(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		format   string
		expected OutputProcessor
		err      string
	}{
		{
			name:     "default output",
			args:     nil,
			expected: &jsonProcessor{},
		},
		{
			name:     "explicit json with options",
			args:     []string{"--output-format", "json", "--json-prefix", "xx"},
			expected: &jsonProcessor{prefix: "xx"},
		},
		{
			name:     "default json with options",
			args:     []string{"--json-indent", "xx"},
			expected: &jsonProcessor{indent: "xx"},
		},
		{
			name:     "raw output",
			args:     []string{"-F", "raw"},
			expected: &rawProcessor{},
		},
		{
			name:     "YAML output",
			args:     []string{"-F", "yaml"},
			expected: &yamlProcessor{},
		},
		{
			name: "template output, no template",
			args: []string{"-F", "template"},
			err:  "Must provide --template or --template-file option",
		},
		{
			name: "unknown output format",
			args: []string{"-F", "oink"},
			err:  "Unrecognized output format 'oink'",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			AddFlags(cmd)
			if err := cmd.ParseFlags(test.args); err != nil {
				t.Fatal(err)
			}
			result, err := SelectOutputProcessor(cmd)
			testy.Error(t, test.err, err)
			if d := diff.Interface(test.expected, result); d != nil {
				t.Error(d)
			}
		})
	}
}

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

func TestRawModeConfig(t *testing.T) {
	cmd := &cobra.Command{}
	mode := &rawMode{}
	mode.config(cmd)

	testOptions(t, []string{}, cmd)
}

func TestRawNew(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		parseErr string
		expected *rawProcessor
		err      string
	}{
		{
			name:     "happy path",
			args:     nil,
			expected: &rawProcessor{},
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
			mode := &rawMode{}
			mode.config(cmd)

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

func TestRawOutput(t *testing.T) {
	input := `{foo bar baz}`
	p := &rawProcessor{}
	t.Run("happy path", func(t *testing.T) {
		buf := &bytes.Buffer{}
		err := p.Output(buf, ioutil.NopCloser(strings.NewReader(input)))
		testy.Error(t, "", err)
		if d := diff.Text(input, buf.String()); d != nil {
			t.Error(d)
		}
	})
	t.Run("write error", func(t *testing.T) {
		err := p.Output(&errWriter{}, ioutil.NopCloser(strings.NewReader(input)))
		testy.Error(t, "errWriter: write error", err)
	})
}

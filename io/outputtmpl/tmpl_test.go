package outputtmpl

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"github.com/flimzy/diff"
	"github.com/flimzy/testy"
	"github.com/go-kivik/kouch"
	"github.com/go-kivik/kouch/internal/test"
	"github.com/go-kivik/kouch/kouchio"
	"github.com/spf13/cobra"
)

func TestTmplModeConfig(t *testing.T) {
	cmd := &cobra.Command{}
	mode := &TmplMode{}
	mode.AddFlags(cmd.PersistentFlags())

	test.Flags(t, []string{"template", "template-file"}, cmd)
}

// func TestTmplNew(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		args     []string
// 		flagsErr string
// 		err      string
// 	}{
// 		{
// 			name: "no options",
// 			err:  "Must provide --template or --template-file option",
// 		},
// 		{
// 			name:     "invalid args",
// 			args:     []string{"--foo"},
// 			flagsErr: "unknown flag: --foo",
// 		},
// 		{
// 			name: "template string & file",
// 			args: []string{"--template", "foo", "--template-file", "bar"},
// 			err:  "Both --template and --template-file specified; must provide only one.",
// 		},
// 		{
// 			name: "invalid template string",
// 			args: []string{"--template", "{{ .foo }"},
// 			err:  `template: :1: unexpected "}" in operand`,
// 		},
// 		{
// 			name: "good template string",
// 			args: []string{"--template", "{{ .foo }}"},
// 		},
// 		{
// 			name: "invalid template file",
// 			args: []string{"--template-file", "./test/template1.html"},
// 			err:  `template: template1.html:1: unexpected "}" in operand`,
// 		},
// 		{
// 			name: "good template string",
// 			args: []string{"--template-file", "./test/template2.html"},
// 		},
// 	}
// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			cmd := &cobra.Command{}
// 			mode := &tmplMode{}
// 			mode.config(cmd.PersistentFlags())
//
// 			err := cmd.ParseFlags(test.args)
// 			testy.Error(t, test.flagsErr, err)
//
// 			result, err := mode.new(cmd, &bytes.Buffer{})
// 			testy.Error(t, test.err, err)
// 			if result.(*tmplProcessor).template == nil {
// 				t.Errorf("Nil template found after instantiation")
// 			}
// 		})
// 	}
// }

func TestTmplOutput(t *testing.T) {
	tests := []struct {
		name             string
		args             []string
		flagsErr, newErr string
		input            string
		expected         string
		err              string
	}{
		{
			name: "happy path",
			args: []string{"--template", `{{ .foo }}`},
			// template: `{{ .foo }}`,
			input:    `{"foo":"bar", "baz":123, "qux": [1,2,3]}`,
			expected: `bar`,
		},
		{
			name:  "invalid JSON input",
			args:  []string{"--template", `{{ .foo }}`},
			input: "oink",
			err:   `invalid character 'o' looking for beginning of value`,
		},
		{
			name:   "no options",
			args:   []string{},
			newErr: "Must provide --template or --template-file option",
		},
		{
			name:     "invalid args",
			args:     []string{"--foo"},
			flagsErr: "unknown flag: --foo",
		},
		{
			name:   "template string & file",
			args:   []string{"--template", "foo", "--template-file", "bar"},
			newErr: "Both --template and --template-file specified; must provide only one.",
		},
		{
			name:   "invalid template string",
			args:   []string{"--template", "{{ .foo }"},
			newErr: `template: :1: unexpected "}" in operand`,
		},
		{
			name:   "invalid template file",
			args:   []string{"--template-file", "./testdata/template1.html"},
			newErr: `template: template1.html:1: unexpected "}" in operand`,
		},
		{
			name:     "good template file",
			args:     []string{"--template-file", "./testdata/template2.html"},
			input:    `{"foo":"bar", "baz":123, "qux": [1,2,3]}`,
			expected: `bar`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			cmd := &cobra.Command{}
			mode := &TmplMode{}
			mode.AddFlags(cmd.PersistentFlags())

			err := cmd.ParseFlags(test.args)
			testy.Error(t, test.flagsErr, err)

			buf := &bytes.Buffer{}
			ctx = kouch.SetFlags(ctx, cmd.Flags())
			p, err := mode.New(ctx, buf)
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

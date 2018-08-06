package io

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/flimzy/diff"
	"github.com/flimzy/testy"
	"github.com/spf13/cobra"
)

func TestAddFlags(t *testing.T) {
	cmd := &cobra.Command{}
	AddFlags(cmd)

	testOptions(t, []string{"force", "json-escape-html", "json-indent", "json-prefix", "output", "output-format", "template", "template-file"}, cmd)
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

func TestSelectOutput(t *testing.T) {
	type soTest struct {
		name         string
		args         []string
		expectedFd   uintptr
		expectedName string
		err          string
		cleanup      func()
	}
	tests := []soTest{
		{
			name:       "default, stdout",
			expectedFd: 1,
		},
		func() soTest {
			f, err := ioutil.TempFile("", "overwrite")
			if err != nil {
				t.Fatal(err)
			}
			f.Close()
			return soTest{
				name:    "overwrite error",
				args:    []string{"--" + FlagOutputFile, f.Name()},
				err:     "^open /tmp/overwrite\\d+: file exists$",
				cleanup: func() { _ = os.Remove(f.Name()) },
			}
		}(),
		func() soTest {
			f, err := ioutil.TempFile("", "overwrite")
			if err != nil {
				t.Fatal(err)
			}
			f.Close()
			return soTest{
				name:         "clobber",
				args:         []string{"--" + FlagOutputFile, f.Name(), "--force"},
				expectedName: f.Name(),
				cleanup:      func() { _ = os.Remove(f.Name()) },
			}
		}(),
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.cleanup != nil {
				defer test.cleanup()
			}
			cmd := &cobra.Command{}
			AddFlags(cmd)
			cmd.ParseFlags(test.args)
			f, err := SelectOutput(cmd)
			testy.ErrorRE(t, test.err, err)
			if file, ok := f.(*os.File); ok {
				if test.expectedFd != 0 {
					if test.expectedFd != file.Fd() {
						t.Errorf("Unexpected FD: Got %d, expected %d", file.Fd(), test.expectedFd)
					}
				}
				if test.expectedName != "" && test.expectedName != file.Name() {
					t.Errorf("Unexpected name: Got %q, expected %q", file.Name(), test.expectedName)
				}
			} else {
				t.Errorf("Unexpected return type: %T", f)
			}
		})
	}
}

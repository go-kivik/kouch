package io

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/flimzy/diff"
	"github.com/flimzy/testy"
	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch"
	"github.com/spf13/cobra"
)

func TestAddFlags(t *testing.T) {
	cmd := &cobra.Command{}
	AddFlags(cmd.PersistentFlags())

	testOptions(t, []string{"data", "force", "json-escape-html", "json-indent", "json-prefix", "output", "output-format", "stderr", "template", "template-file"}, cmd)
}

func TestSelectOutputProcessor(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		format   string
		expected kouch.OutputProcessor
		err      string
	}{
		{
			name:     "default output",
			args:     nil,
			expected: &errWrapper{&jsonProcessor{}},
		},
		{
			name:     "explicit json with options",
			args:     []string{"--output-format", "json", "--json-prefix", "xx"},
			expected: &errWrapper{&jsonProcessor{prefix: "xx"}},
		},
		{
			name:     "default json with options",
			args:     []string{"--json-indent", "xx"},
			expected: &errWrapper{&jsonProcessor{indent: "xx"}},
		},
		{
			name:     "raw output",
			args:     []string{"-F", "raw"},
			expected: &errWrapper{&rawProcessor{}},
		},
		{
			name:     "YAML output",
			args:     []string{"-F", "yaml"},
			expected: &errWrapper{&yamlProcessor{}},
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
			AddFlags(cmd.PersistentFlags())
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
		{
			name: "Missing parent dir",
			args: []string{"--" + FlagOutputFile, "./foo/bar/baz"},
			err:  "open ./foo/bar/baz: no such file or directory",
		},
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
			AddFlags(cmd.PersistentFlags())
			cmd.ParseFlags(test.args)
			f, err := SelectOutput(cmd)
			if err != nil {
				t.Fatal(err)
			}
			switch file := f.(type) {
			case *os.File:
				if test.expectedFd != 0 {
					if test.expectedFd != file.Fd() {
						t.Errorf("Unexpected FD: Got %d, expected %d", file.Fd(), test.expectedFd)
					}
				}
				if test.expectedName != "" && test.expectedName != file.Name() {
					t.Errorf("Unexpected name: Got %q, expected %q", file.Name(), test.expectedName)
				}
			case *delayedOpenWriter:
				if test.expectedName != "" && test.expectedName != file.filename {
					t.Errorf("Unexpected name: Got %q, expected %q", file.filename, test.expectedName)
				}
			default:
				t.Errorf("Unexpected return type: %T", f)
			}

			_, err = f.Write([]byte("foo"))
			testy.ErrorRE(t, test.err, err)
		})
	}
}

func TestRedirStderr(t *testing.T) {
	var stderr = os.Stderr
	defer func() {
		// Restore original setting
		os.Stderr = stderr
	}()
	type rsTest struct {
		name     string
		args     []string
		expected string
		err      string
		status   int
		cleanup  func()
	}
	tests := []rsTest{
		{
			name:     "No redirection",
			args:     nil,
			expected: "/dev/stderr",
		},
		{
			name:   "Dir doesn't exist",
			args:   []string{"--stderr", "./does_not_exist/foo"},
			err:    "open ./does_not_exist/foo: no such file or directory",
			status: chttp.ExitWriteError,
		},
		func() rsTest {
			tmpDir, err := ioutil.TempDir("", "stderrRedir-")
			if err != nil {
				t.Fatal(err)
			}
			return rsTest{
				name:     "redir to file",
				args:     []string{"--stderr", tmpDir + "/foo"},
				expected: tmpDir + "/foo",
				cleanup:  func() { _ = os.RemoveAll(tmpDir) },
			}
		}(),
		func() rsTest {
			f, err := ioutil.TempFile("", "stderrRedir-")
			if err != nil {
				t.Fatal(err)
			}
			tmpfile := f.Name()
			_ = f.Close()
			return rsTest{
				name:    "file already exists",
				args:    []string{"--stderr", tmpfile},
				err:     "open /tmp/stderrRedir-\\d+: file exists",
				status:  chttp.ExitWriteError,
				cleanup: func() { _ = os.Remove(tmpfile) },
			}
		}(),
		func() rsTest {
			f, err := ioutil.TempFile("", "stderrRedir-")
			if err != nil {
				t.Fatal(err)
			}
			tmpfile := f.Name()
			_ = f.Close()
			return rsTest{
				name:     "file already exists, + clobber",
				args:     []string{"--stderr", tmpfile, "--force"},
				expected: tmpfile,
				cleanup:  func() { _ = os.Remove(tmpfile) },
			}
		}(),
		{
			name:     "stdout",
			args:     []string{"--stderr", "-"},
			expected: "/dev/stdout",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.cleanup != nil {
				defer test.cleanup()
			}
			t.Run("group", func(t *testing.T) {
				cmd := &cobra.Command{}
				AddFlags(cmd.PersistentFlags())
				cmd.ParseFlags(test.args)
				err := RedirStderr(cmd.Flags())
				testy.ExitStatusErrorRE(t, test.err, test.status, err)
				filename := os.Stderr.Name()
				if filename != test.expected {
					t.Errorf("Unexpected filename: %s", filename)
				}
			})
		})
	}
}

func TestSelectInput(t *testing.T) {
	type siTest struct {
		name     string
		args     []string
		err      string
		status   int
		expected string
		cleanup  func()
	}
	tests := []siTest{
		func() siTest {
			stdin := os.Stdin
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatal(err)
			}
			os.Stdin = r
			go func() {
				w.Write([]byte("stdin data"))
				w.Close()
			}()
			return siTest{
				name:     "defaults",
				expected: "stdin data",
				cleanup:  func() { os.Stdin = stdin },
			}
		}(),
		{
			name:     "input string",
			args:     []string{"--" + kouch.FlagData, "some data"},
			expected: "some data",
		},
		func() siTest {
			f, err := ioutil.TempFile("", "overwrite")
			if err != nil {
				t.Fatal(err)
			}
			if _, e := f.Write([]byte("file data")); e != nil {
				t.Fatal(e)
			}
			f.Close()
			return siTest{
				name:     "read from file",
				args:     []string{"--" + kouch.FlagData, "@" + f.Name()},
				expected: "file data",
				cleanup:  func() { _ = os.Remove(f.Name()) },
			}
		}(),
		{
			name:   "read from missing file",
			args:   []string{"--" + kouch.FlagData, "@missingfile.txt"},
			err:    "open missingfile.txt: no such file or directory",
			status: chttp.ExitReadError,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.cleanup != nil {
				defer test.cleanup()
			}
			cmd := &cobra.Command{}
			AddFlags(cmd.PersistentFlags())
			cmd.ParseFlags(test.args)
			f, err := SelectInput(cmd)
			testy.ExitStatusError(t, test.err, test.status, err)
			content, err := ioutil.ReadAll(f)
			if err != nil {
				t.Fatal(err)
			}
			if d := diff.Text(test.expected, content); d != nil {
				t.Error(d)
			}
		})
	}
}

package io

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/flimzy/diff"
	"github.com/flimzy/testy"
	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func TestAddFlags(t *testing.T) {
	cmd := &cobra.Command{}
	AddFlags(cmd.PersistentFlags())

	testOptions(t, []string{"data", "data-json", "data-yaml", "force", "json-escape-html", "json-indent", "json-prefix", "output", "output-format", "stderr", "template", "template-file"}, cmd)
}

func TestSelectOutputProcessor(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		format   string
		expected io.Writer
		err      string
	}{
		// {
		// 	name:     "default output",
		// 	args:     nil,
		// 	expected: &errWrapper{&jsonProcessor{}},
		// },
		// {
		// 	name:     "explicit json with options",
		// 	args:     []string{"--output-format", "json", "--json-prefix", "xx"},
		// 	expected: &errWrapper{&jsonProcessor{prefix: "xx"}},
		// },
		// {
		// 	name:     "default json with options",
		// 	args:     []string{"--json-indent", "xx"},
		// 	expected: &errWrapper{&jsonProcessor{indent: "xx"}},
		// },
		{
			name:     "raw output",
			args:     []string{"-F", "raw"},
			expected: &exitStatusWriter{&nopCloser{&bytes.Buffer{}}},
		},
		// {
		// 	name:     "YAML output",
		// 	args:     []string{"-F", "yaml"},
		// 	expected: &errWrapper{&yamlProcessor{}},
		// },
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
			result, err := selectOutputProcessor(cmd, &bytes.Buffer{})
			testy.Error(t, test.err, err)
			if d := diff.Interface(test.expected, result); d != nil {
				t.Error(d)
			}
		})
	}
}

func TestOpen(t *testing.T) {
	newFlags := func() *pflag.FlagSet {
		f := pflag.NewFlagSet("foo", 1)
		f.String(kouch.FlagOutputFile, "", "x")
		f.Bool(kouch.FlagClobber, false, "x")
		return f
	}
	type oTest struct {
		flags        *pflag.FlagSet
		flagName     string
		expectedNil  bool
		expectedFd   uintptr
		expectedName string
		err          string
		writeErr     string
	}
	tests := testy.NewTable()
	tests.Add("no flag defined", oTest{
		flags:    newFlags(),
		flagName: "foo",
		err:      "flag accessed but not defined: foo",
	})
	tests.Add("default", oTest{
		flags:       newFlags(),
		flagName:    kouch.FlagOutputFile,
		expectedNil: true,
	})
	tests.Add("stdout", oTest{
		flags: func() *pflag.FlagSet {
			f := newFlags()
			f.Set(kouch.FlagOutputFile, "-")
			return f
		}(),
		flagName:     kouch.FlagOutputFile,
		expectedFd:   1,
		expectedName: "/dev/stdout",
	})
	tests.Add("stderr", oTest{
		flags: func() *pflag.FlagSet {
			f := newFlags()
			f.Set(kouch.FlagOutputFile, "%")
			return f
		}(),
		flagName:     kouch.FlagOutputFile,
		expectedFd:   2,
		expectedName: "/dev/stderr",
	})
	tests.Add("overwrite error", func(t *testing.T) interface{} {
		file, err := ioutil.TempFile("", "overwrite")
		if err != nil {
			t.Fatal(err)
		}
		tests.Cleanup(func() error {
			return os.Remove(file.Name())
		})
		file.Close()

		flags := newFlags()
		flags.Set(kouch.FlagOutputFile, file.Name())
		return oTest{
			flags:        flags,
			flagName:     kouch.FlagOutputFile,
			expectedName: file.Name(),
			writeErr:     "^open /tmp/overwrite\\d+: file exists$",
		}
	})
	tests.Add("missing parent dir", oTest{
		flags: func() *pflag.FlagSet {
			f := newFlags()
			f.Set(kouch.FlagOutputFile, "./foo/bar/baz")
			return f
		}(),
		flagName:     kouch.FlagOutputFile,
		expectedName: "./foo/bar/baz",
		writeErr:     "open ./foo/bar/baz: no such file or directory",
	})
	tests.Add("clobber", func(t *testing.T) interface{} {
		file, err := ioutil.TempFile("", "overwrite")
		if err != nil {
			t.Fatal(err)
		}
		file.Close()

		flags := newFlags()
		flags.Set(kouch.FlagOutputFile, file.Name())
		flags.Set(kouch.FlagClobber, "true")

		tests.Cleanup(func() error {
			return os.Remove(file.Name())
		})

		return oTest{
			flags:        flags,
			flagName:     kouch.FlagOutputFile,
			expectedName: file.Name(),
		}
	})

	tests.Run(t, func(t *testing.T, test oTest) {
		f, err := open(test.flags, test.flagName)
		testy.Error(t, test.err, err)
		if test.expectedNil {
			if f != nil {
				t.Errorf("Expected nil, got %T", f)
			}
			return
		}

		testFile(t, f, test.expectedFd, test.expectedName)

		_, err = f.Write([]byte("foo"))
		testy.ErrorRE(t, test.writeErr, err)
	})
}

func TestSetOutput(t *testing.T) {
	type soTest struct {
		name       string
		args       []string
		outputFd   uintptr
		outputName string
		stderrFd   uintptr
		stderrName string
		headFd     uintptr
		headName   string
		err        string
		cleanup    func()
	}
	tests := []soTest{
		{
			name:       "default, stdout",
			outputFd:   1,
			outputName: "/dev/stdout",
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
			ctx := context.Background()
			var err error
			ctx, err = setOutput(ctx, cmd.Flags())
			if err != nil {
				t.Fatal(err)
			}
			f := kouch.Output(ctx)
			testFile(t, f, test.outputFd, test.outputName)

			_, err = f.Write([]byte("foo"))
			testy.ErrorRE(t, test.err, err)
		})
	}
}

func testFile(t *testing.T, f io.Writer, expectedFd uintptr, expectedName string) {
	switch file := f.(type) {
	case *os.File:
		if expectedFd != 0 {
			if expectedFd != file.Fd() {
				t.Errorf("Unexpected FD: Got %d, expected %d", file.Fd(), expectedFd)
			}
		}
		if expectedName != file.Name() {
			t.Errorf("Unexpected name: Got %q, expected %q", file.Name(), expectedName)
		}
	case *delayedOpenWriter:
		if expectedName != file.filename {
			t.Errorf("Unexpected name: Got %q, expected %q", file.filename, expectedName)
		}
	default:
		t.Errorf("Unexpected return type: %T", f)
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
		{
			name:   "too much data",
			args:   []string{"--" + kouch.FlagData, "foo", "--" + kouch.FlagDataJSON, "bar"},
			err:    "Only one data option may be provided",
			status: chttp.ExitFailedToInitialize,
		},
		{
			name:   "invalid json input",
			args:   []string{"--" + kouch.FlagDataJSON, "invalid"},
			err:    "invalid character 'i' looking for beginning of value",
			status: chttp.ExitPostError,
		},
		{
			name:     "json input",
			args:     []string{"--" + kouch.FlagDataJSON, `{ "_id": "foo" }`},
			expected: `{"_id":"foo"}`,
		},
		{
			name:   "invalid yaml input",
			args:   []string{"--" + kouch.FlagDataYAML, `{]}`},
			err:    "yaml: did not find expected node content",
			status: chttp.ExitPostError,
		},
		{
			name:     "yaml input",
			args:     []string{"--" + kouch.FlagDataYAML, `_id: foo`},
			expected: `{"_id":"foo"}`,
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

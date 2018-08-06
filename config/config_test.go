package config

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/flimzy/diff"
	"github.com/flimzy/testy"
	"github.com/go-kivik/kouch"
	"github.com/spf13/cobra"
)

var expectedConf = &kouch.Config{DefaultContext: "foo",
	Contexts: []kouch.NamedContext{
		{
			Name:    "foo",
			Context: &kouch.Context{Root: "http://foo.com/"},
		},
	},
}

func TestReadConfigFile(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *kouch.Config
		err      string
	}{
		{
			name: "not found",
			err:  "^open /tmp/TestReadConfigFile_not_found-\\d+/config: no such file or directory$",
		},
		{
			name: "yaml input",
			input: `default-context: foo
contexts:
- context:
    root: http://foo.com/
  name: foo
`,
			expected: expectedConf,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tmpDir := new(string)
			defer testy.TempDir(t, tmpDir)()
			file := path.Join(*tmpDir, "config")
			if test.input != "" {
				if err := ioutil.WriteFile(file, []byte(test.input), 0777); err != nil {
					t.Fatal(err)
				}
			}
			conf, err := readConfigFile(file)
			testy.ErrorRE(t, test.err, err)
			if d := diff.Interface(test.expected, conf); d != nil {
				t.Fatal(d)
			}
		})
	}
}

func TestReadConfig(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		env      map[string]string
		args     []string
		expected *kouch.Config
		err      string
	}{
		{
			name:     "no config",
			expected: &kouch.Config{},
		},
		{
			name: "default config only",
			files: map[string]string{
				".kouch/config": `default-context: foo
contexts:
- context:
    root: http://foo.com/
  name: foo
`,
			},
			expected: expectedConf,
		},
		{
			name: "specific config file",
			files: map[string]string{
				"kouch.yaml": `default-context: foo
contexts:
- context:
    root: http://foo.com/
  name: foo
`,
				".kouch/config": `default-context: bar
contexts:
- context:
    root: http://bar.com/
  name: bar
`,
			},
			args:     []string{"--kouchconfig", "${HOME}/kouch.yaml"},
			expected: expectedConf,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tmpDir := new(string)
			defer testy.TempDir(t, tmpDir)()
			defer testy.RestoreEnv()()
			env := map[string]string{"HOME": *tmpDir}
			for k, v := range test.env {
				env[k] = strings.Replace(v, "${HOME}", *tmpDir, -1)
			}
			testy.SetEnv(env)
			for filename, content := range test.files {
				file := path.Join(*tmpDir, filename)
				if err := os.MkdirAll(path.Dir(file), 0777); err != nil {
					t.Fatal(err)
				}
				if err := ioutil.WriteFile(file, []byte(content), 0777); err != nil {
					t.Fatal(err)
				}
			}

			cmd := &cobra.Command{}
			AddFlags(cmd)
			for i, v := range test.args {
				test.args[i] = strings.Replace(v, "${HOME}", *tmpDir, -1)
			}
			cmd.ParseFlags(test.args)

			conf, err := ReadConfig(cmd)
			testy.ErrorRE(t, test.err, err)
			if d := diff.Interface(test.expected, conf); d != nil {
				t.Fatal(d)
			}
		})
	}
}

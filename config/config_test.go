package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/flimzy/diff"
	"github.com/flimzy/testy"
	"github.com/go-kivik/kouch"
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
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tmpDir := new(string)
			defer testy.TempDir(t, tmpDir)()
			defer testy.RestoreEnv()()
			env := map[string]string{"HOME": *tmpDir}
			for k, v := range test.env {
				env[k] = v
			}
			testy.SetEnv(env)
			for filename, content := range test.files {
				file := path.Join(*tmpDir, filename)
				fmt.Printf("gonna create %s\n", file)
				fmt.Printf("Creating dir: %s\n", path.Dir(file))
				if err := os.MkdirAll(path.Dir(file), 0777); err != nil {
					t.Fatal(err)
				}
				if err := ioutil.WriteFile(file, []byte(content), 0777); err != nil {
					t.Fatal(err)
				}
			}
			conf, err := ReadConfig()
			testy.ErrorRE(t, test.err, err)
			if d := diff.Interface(test.expected, conf); d != nil {
				t.Fatal(d)
			}
		})
	}
}

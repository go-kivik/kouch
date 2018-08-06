package config

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/flimzy/diff"
	"github.com/flimzy/testy"
	"github.com/go-kivik/kouch"
)

func TestReadConfigFile(t *testing.T) {
	expectedConf := &kouch.Config{DefaultContext: "foo",
		Contexts: []kouch.NamedContext{
			{
				Name:    "foo",
				Context: &kouch.Context{Root: "http://foo.com/"},
			},
		},
	}
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

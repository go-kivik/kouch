package config

import (
	"testing"

	"github.com/flimzy/testy"
	"github.com/go-kivik/kouch/internal/test"

	_ "github.com/go-kivik/kouch/cmd/kouch/root"
)

func TestConfigCmd(t *testing.T) {
	tests := testy.NewTable()
	tests.Add("no config", test.CmdTest{
		Stdout: "{}",
	})
	tests.Add("from command line", test.CmdTest{
		Args: []string{"--root", "foo.com", "-F", "yaml"},
		Stdout: `contexts:
- context:
    root: foo.com
  name: $dynamic$
default-context: $dynamic$
`,
	})

	tests.Run(t, test.ValidateCmdTest([]string{"config", "view"}))
}

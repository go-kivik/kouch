package config

import (
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/go-kivik/kouch"
	"github.com/go-kivik/kouch/cmd/kouch/registry"
	"github.com/spf13/cobra"
)

const (
	FlagKouchConfigFile = "kouchconfig"
	envKouchConfigFiles = "KOUCHCONFIG"
	kouchHome           = ".kouch"
)

func init() {
	registry.Register([]string{}, func(_ *kouch.Context) *cobra.Command {
		cmd := &cobra.Command{
			Use:   "config",
			Short: "Modify kouchconfig files",
			Long: `Modify kouchconfig files using subcommands.

The loading order follows these rules:

  1. If the --` + FlagKouchConfigFile + ` flag is set, then only that file is loaded.  The flag may only be set once and no merging takes place.
  2. If the $` + envKouchConfigFiles + ` environment variable is set, then it is used as a list of paths (normal path delimitting rules for your system).  These paths are merged.  When a value is modified, it is modified in the file that defines the stanza.  When a value is created, it is created in the first file that exists.  If no files in the chain exist, then it creates the last file in the list.
  3. Otherwise, ` + path.Join("${HOME}", kouchHome) + `/config is used and no merging takes place.`,
		}
		cmd.Flags().String(FlagKouchConfigFile, "", "Path to the kouchconfig file to use for CLI requests.")
		return cmd
	})

	registry.Register([]string{"config"}, func(cx *kouch.Context) *cobra.Command {
		cmd := &cobra.Command{
			Use:   "view",
			Short: "Display merged kouchconfig settings or a specified kouchconfig file",
			RunE:  viewConfig(cx),
		}
		return cmd
	})
}

func viewConfig(cx *kouch.Context) func(*cobra.Command, []string) error {
	return func(_ *cobra.Command, _ []string) error {
		return cx.Outputer.Output(os.Stdout, ioutil.NopCloser(strings.NewReader(`{"foo":"bar"}`)))
	}
}

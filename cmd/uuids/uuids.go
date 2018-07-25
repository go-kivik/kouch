package get

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/go-kivik/kouch/cmd/registry"
	"github.com/go-kivik/kouch/log"
)

func init() {
	registry.Register([]string{"get"}, func(log log.Logger, conf *viper.Viper) *cobra.Command {
		cmd := &cobra.Command{
			Use:   "uuids",
			Short: "Returns one or more server-generated UUIDs",
			Long: `Returns one or more Universally Unique Identifiers (UUIDs) from the
CouchDB server.`,
			Run: uuidsCmd,
		}
		cmd.Flags().IntP("count", "C", 1, "Number of UUIDs to return")
		return cmd
	})
}

func uuidsCmd(cmd *cobra.Command, _ []string) {
	count, err := cmd.Flags().GetInt("count")
	if err != nil {
		panic(err.Error())
	}
	url, err := cmd.Flags().GetString("url")
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("%d UUIDs coming right up, from %s\n", count, url)
}

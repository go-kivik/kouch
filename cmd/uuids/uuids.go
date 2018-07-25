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
		var count int
		cmd := &cobra.Command{
			Use:   "uuids",
			Short: "Returns one or more server-generated UUIDs",
			Long: `Returns one or more Universally Unique Identifiers (UUIDs) from the
CouchDB server.`,
			Run: uuidsCmd(&count),
		}
		cmd.Flags().IntVarP(&count, "count", "C", 1, "Number of UUIDs to return")
		return cmd
	})
}

func uuidsCmd(count *int) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		fmt.Printf("%d UUIDs coming right up\n", *count)
	}
}

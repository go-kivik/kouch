package get

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/go-kivik/kouch/cmd/registry"
	"github.com/go-kivik/kouch/log"
)

func init() {
	registry.Register(nil, func(log log.Logger, conf *viper.Viper) *cobra.Command {
		return &cobra.Command{
			Use:   "get",
			Short: "Display one or more resources.",
		}
	})
}

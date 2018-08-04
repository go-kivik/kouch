package get

import (
	"github.com/spf13/cobra"

	"github.com/go-kivik/kouch"
	"github.com/go-kivik/kouch/cmd/kouch/registry"
)

func init() {
	registry.Register(nil, func(cx *kouch.Context) *cobra.Command {
		return &cobra.Command{
			Use:   "get",
			Short: "Display one or more resources.",
		}
	})
}

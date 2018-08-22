package get

import (
	"github.com/spf13/cobra"

	"github.com/go-kivik/kouch/cmd/kouch/registry"
)

func init() {
	registry.Register(nil, getCmd)
}

func getCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Display one or more resources.",
	}
}

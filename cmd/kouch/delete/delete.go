package delete

import (
	"github.com/go-kivik/kouch/cmd/kouch/registry"
	"github.com/spf13/cobra"
)

func init() {
	registry.Register(nil, deleteCmd)
}

func deleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "delete",
		Aliases: []string{"del"},
		Short:   "Delete a resource.",
	}
}

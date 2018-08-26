package create

import (
	"github.com/go-kivik/kouch/cmd/kouch/registry"
	"github.com/spf13/cobra"
)

func init() {
	registry.Register(nil, createCmd)
}

func createCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create",
		Short: "Create a resource.",
	}
}

package get

import (
	"github.com/spf13/cobra"

	"github.com/go-kivik/kouch/cmd/kouch/registry"
)

func init() {
	registry.Register(nil, &cobra.Command{
		Use:   "put",
		Short: "Create or update an existing resource.",
	})
}

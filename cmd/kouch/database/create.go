package database

import (
	"github.com/go-kivik/kouch/cmd/kouch/registry"
	"github.com/go-kivik/kouch/target"
	"github.com/spf13/cobra"
)

func init() {
	registry.Register([]string{"create"}, createDbCmd)
}

func createDbCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "database [target]",
		Aliases: []string{"db"},
		Short:   "Creates a new database.",
		Long: "Creates a new database.\n\n" +
			target.HelpText(target.Database),
		RunE: createDatabaseCmd,
	}
	return cmd
}

func createDatabaseCmd(cmd *cobra.Command, _ []string) error {
	return nil
}

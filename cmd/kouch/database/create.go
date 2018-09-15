package database

import (
	"context"
	"net/http"

	"github.com/go-kivik/kouch"
	"github.com/go-kivik/kouch/cmd/kouch/registry"
	"github.com/go-kivik/kouch/internal/util"
	"github.com/go-kivik/kouch/target"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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
			target.HelpText(kouch.TargetDatabase),
		RunE: createDatabaseCmd,
	}
	cmd.Flags().IntP(kouch.FlagShards, kouch.FlagShortShards, 0, "Shards, aka the number of range partitions.")
	return cmd
}

func createDatabaseCmd(cmd *cobra.Command, _ []string) error {
	ctx := kouch.GetContext(cmd)
	o, err := createDatabaseOpts(ctx, cmd.Flags())
	if err != nil {
		return err
	}
	if err := validateTarget(o.Target); err != nil {
		return err
	}
	return util.ChttpDo(ctx, http.MethodPut, util.DatabasePath(o), o)
}

func createDatabaseOpts(ctx context.Context, flags *pflag.FlagSet) (*kouch.Options, error) {
	o, err := util.CommonOptions(ctx, kouch.TargetDatabase, flags)

	if e := o.SetParamInt(flags, kouch.FlagShards); e != nil {
		return nil, e
	}
	return o, err
}

func validateTarget(target *kouch.Target) error {
	return nil
}

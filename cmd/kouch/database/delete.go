package database

import (
	"context"
	"net/http"

	"github.com/go-kivik/kouch"
	"github.com/go-kivik/kouch/cmd/kouch/registry"
	"github.com/go-kivik/kouch/internal/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func init() {
	registry.Register([]string{"delete"}, deleteDbCmd)
}

func deleteDbCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "database [target]",
		Aliases: []string{"db"},
		Short:   "Deletes a database.",
		Long: "Deletes a database.\n\n" +
			kouch.TargetHelpText(kouch.TargetDatabase),
		RunE: deleteDatabaseCmd,
	}
	return cmd
}

func deleteDatabaseCmd(cmd *cobra.Command, _ []string) error {
	ctx := kouch.GetContext(cmd)
	o, err := deleteDatabaseOpts(ctx, cmd.Flags())
	if err != nil {
		return err
	}
	return util.ChttpDo(ctx, http.MethodDelete, util.DatabasePath(o), o)
}

func deleteDatabaseOpts(ctx context.Context, flags *pflag.FlagSet) (*kouch.Options, error) {
	o, err := util.CommonOptions(ctx, kouch.TargetDatabase, flags)

	if e := o.SetParamInt(flags, kouch.FlagShards); e != nil {
		return nil, e
	}
	return o, err
}

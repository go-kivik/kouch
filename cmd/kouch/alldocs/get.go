package alldocs

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
	registry.Register([]string{"get"}, getAllDocsCmd)
}

func getAllDocsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "alldocs [target]",
		Aliases: []string{"doc"},
		Short:   "Fetches all documents in a database.",
		Long: "Fetches all documents in a database, subject to possible restrictions.\n\n" +
			kouch.TargetHelpText(kouch.TargetDatabase),
		RunE: getAllDocumentsCmd,
	}
	return cmd
}

func getAllDocumentsCmd(cmd *cobra.Command, _ []string) error {
	ctx := kouch.GetContext(cmd)
	o, err := getAllDocsOpts(ctx, cmd.Flags())
	if err != nil {
		return err
	}
	return getAllDocs(ctx, o)
}

func getAllDocsOpts(ctx context.Context, flags *pflag.FlagSet) (*kouch.Options, error) {
	o, err := util.CommonOptions(ctx, kouch.TargetDocument, flags)
	if err != nil {
		return nil, err
	}

	return o, nil
}

func getAllDocs(ctx context.Context, o *kouch.Options) error {
	if err := validateTarget(o.Target); err != nil {
		return err
	}
	return util.ChttpDo(ctx, http.MethodGet, util.DocPath(o), o)
}

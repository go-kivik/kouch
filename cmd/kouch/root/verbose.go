package root

import (
	"context"

	"github.com/go-kivik/kouch"
	"github.com/spf13/cobra"
)

func verbose(ctx context.Context, cmd *cobra.Command) (context.Context, error) {
	verbose, err := cmd.Flags().GetBool(flagVerbose)
	if err != nil {
		return ctx, err
	}
	if !verbose {
		return ctx, nil
	}
	ctx = kouch.SetVerbose(ctx, true)
	return ctx, nil
}

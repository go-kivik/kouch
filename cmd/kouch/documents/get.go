package documents

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch"
	"github.com/go-kivik/kouch/cmd/kouch/registry"
	"github.com/go-kivik/kouch/target"
	"github.com/spf13/cobra"
)

func init() {
	registry.Register([]string{"get"}, getDocCmd())
}

func getDocCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "document [target]",
		Aliases: []string{"doc"},
		Short:   "Fetches a single document.",
		Long: "Fetches a single document.\n\n" +
			target.HelpText(target.Document),
		RunE: getDocumentCmd,
	}
	f := cmd.Flags()
	f.String(kouch.FlagDocument, "", "The document ID. May be provided with the target in the format {id}.")
	f.String(kouch.FlagDatabase, "", "The database. May be provided with the target in the format /{db}/{id}.")
	f.StringP(kouch.FlagRev, kouch.FlagShortRev, "", "Retrieves document of specified revision.")
	f.String(kouch.FlagIfNoneMatch, "", "Optionally fetch the document, only if the current rev does not match the one provided")

	cmd.PersistentFlags().BoolP(kouch.FlagHead, kouch.FlagShortHead, false, "Fetch the headers only.")
	f.Bool(flagIncludeAttachments, false, "Include attachments bodies in response.")
	f.Bool(flagIncludeAttEncoding, false, "Include encoding information in attachment stubs for compressed attachments.")
	f.StringSlice(flagAttsSince, nil, "Include attachments only since, but not including, the specified revisions.")
	f.Bool(flagIncludeConflicts, false, "Include document conflicts information.")
	f.Bool(flagIncludeDeletedConflicts, false, "Include information about deleted conflicted revisions.")
	f.Bool(flagForceLatest, false, `Force retrieving latest “leaf” revision, no matter what rev was requested.`)
	f.Bool(flagIncludeLocalSeq, false, "Include last update sequence for the document.")
	f.Bool(flagMeta, false, "Same as: --"+flagIncludeConflicts+" --"+flagIncludeDeletedConflicts+" --"+flagRevsInfo)
	f.StringSlice(flagOpenRevs, nil, "Retrieve documents of specified leaf revisions. May use the value 'all' to return all leaf revisions.")
	f.Bool(flagRevs, false, "Include list of all known document revisions.")
	f.Bool(flagRevsInfo, false, "Include detailed information for all known document revisions")
	return cmd
}

func getDocumentCmd(cmd *cobra.Command, args []string) error {
	ctx := kouch.GetContext(cmd)
	o, err := getDocumentOpts(cmd, args)
	if err != nil {
		return err
	}
	result, err := getDocument(o)
	if err != nil {
		return err
	}
	if o.Head {
		defer result.Close()
		_, err = io.Copy(kouch.Output(ctx), result)
		return err
	}
	return kouch.Outputer(ctx).Output(kouch.Output(ctx), result)
}

func getDocumentOpts(cmd *cobra.Command, _ []string) (*kouch.Options, error) {
	ctx := kouch.GetContext(cmd)
	o := kouch.NewOptions()
	if tgt := kouch.GetTarget(ctx); tgt != "" {
		var err error
		o.Target, err = target.Parse(target.Document, tgt)
		if err != nil {
			return nil, err
		}
	}

	if err := o.Target.DocumentFromFlags(cmd.Flags()); err != nil {
		return nil, err
	}
	if err := o.Target.DatabaseFromFlags(cmd.Flags()); err != nil {
		return nil, err
	}
	if err := o.SetHead(cmd.Flags()); err != nil {
		return nil, err
	}

	if defCtx, err := kouch.Conf(ctx).DefaultCtx(); err == nil {
		if o.Root == "" {
			o.Root = defCtx.Root
		}
	}
	var err error
	o.Options.IfNoneMatch, err = cmd.Flags().GetString(kouch.FlagIfNoneMatch)
	if err != nil {
		return nil, err
	}
	if e := o.SetParamString(cmd.Flags(), kouch.FlagRev); e != nil {
		return nil, e
	}
	for _, flag := range []string{flagAttsSince, flagOpenRevs} {
		if e := o.SetParamStringSlice(cmd.Flags(), flag); e != nil {
			return nil, e
		}
	}

	for _, flag := range []string{
		flagIncludeAttachments, flagIncludeAttEncoding, flagIncludeConflicts,
		flagIncludeDeletedConflicts, flagForceLatest, flagIncludeLocalSeq,
		flagMeta, flagRevs, flagRevsInfo,
	} {
		if e := o.SetParamBool(cmd.Flags(), flag); e != nil {
			return nil, e
		}
	}

	return o, nil
}

func getDocument(o *kouch.Options) (io.ReadCloser, error) {
	if err := validateTarget(o.Target); err != nil {
		return nil, err
	}
	c, err := chttp.New(context.TODO(), o.Root)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/%s/%s", url.QueryEscape(o.Database), chttp.EncodeDocID(o.Document))
	method := http.MethodGet
	if o.Head {
		method = http.MethodHead
	}
	res, err := c.DoReq(context.TODO(), method, path, o.Options)
	if err != nil {
		return nil, err
	}
	if err = chttp.ResponseError(res); err != nil {
		return nil, err
	}
	if o.Head {
		_ = res.Body.Close()
		r, w := io.Pipe()
		go func() {
			err := res.Header.Write(w)
			w.CloseWithError(err)
		}()
		return r, nil
	}
	return res.Body, nil
}

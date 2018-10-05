package documents

import (
	"strings"

	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch"
	"github.com/go-kivik/kouch/internal/errors"
	"github.com/spf13/pflag"
)

func param(flagName string) string {
	return strings.Replace(flagName, "-", "_", -1)
}

func setBatch(o *kouch.Options, f *pflag.FlagSet) error {
	v, err := f.GetBool(kouch.FlagBatch)
	if err == nil && v {
		o.Query().Add(param(kouch.FlagBatch), "ok")
	}
	return err
}

func validateTarget(t *kouch.Target) error {
	if t.Filename != "" {
		panic("non-nil filename")
	}
	if t.Document == "" {
		return errors.NewExitError(chttp.ExitFailedToInitialize, "No document ID provided")
	}
	if t.Database == "" {
		return errors.NewExitError(chttp.ExitFailedToInitialize, "No database name provided")
	}
	if t.Root == "" {
		return errors.NewExitError(chttp.ExitFailedToInitialize, "No root URL provided")
	}
	return nil
}

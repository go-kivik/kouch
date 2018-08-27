package kouch

import (
	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch/internal/errors"
	"github.com/spf13/pflag"
)

// Target is a parsed target passed on the command line.
type Target struct {
	// Root is the root URL.
	Root string
	// Database is the database name.
	Database string
	// DocID is the document ID.
	Document string
	// Filename is the attachment filename.
	Filename string
	// Username is the HTTP Basic Auth username
	Username string
	// Password is the HTTP Basic Auth password
	Password string
}

// DocumentFromFlags sets t.DocID from the passed flagset.
func (t *Target) DocumentFromFlags(flags *pflag.FlagSet) error {
	if flag := flags.Lookup(FlagDocument); flag == nil {
		return nil
	}
	id, err := flags.GetString(FlagDocument)
	if err != nil {
		return err
	}
	if id == "" {
		return nil
	}
	if t.Document != "" {
		return errors.NewExitError(chttp.ExitFailedToInitialize,
			"Must not use --%s and pass document ID as part of the target", FlagDocument)
	}
	t.Document = id
	return nil
}

// DatabaseFromFlags sets t.Database from the passed flagset.
func (t *Target) DatabaseFromFlags(flags *pflag.FlagSet) error {
	if flag := flags.Lookup(FlagDatabase); flag == nil {
		return nil
	}
	db, err := flags.GetString(FlagDatabase)
	if err != nil {
		return err
	}
	if db == "" {
		return nil
	}
	if t.Database != "" {
		return errors.NewExitError(chttp.ExitFailedToInitialize,
			"Must not use --%s and pass database as part of the target", FlagDatabase)
	}
	t.Database = db
	return nil
}

// FilenameFromFlags sets t.Filename from the passed flagset.
func (t *Target) FilenameFromFlags(flags *pflag.FlagSet) error {
	if flag := flags.Lookup(FlagFilename); flag == nil {
		return nil
	}
	fn, err := flags.GetString(FlagFilename)
	if err != nil {
		return err
	}
	if fn == "" {
		return nil
	}
	if t.Filename != "" {
		return errors.NewExitError(chttp.ExitFailedToInitialize,
			"Must not use --%s and pass separate filename", FlagFilename)
	}
	t.Filename = fn
	return nil
}

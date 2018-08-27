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

var duplicateConfigErrors = map[string]error{
	FlagDatabase: errors.NewExitError(chttp.ExitFailedToInitialize,
		"Must not use --%s and pass database as part of the target", FlagDatabase),
	FlagDocument: errors.NewExitError(chttp.ExitFailedToInitialize,
		"Must not use --%s and pass document ID as part of the target", FlagDocument),
	FlagFilename: errors.NewExitError(chttp.ExitFailedToInitialize,
		"Must not use --%s and pass separate filename", FlagFilename),
}

func setFromFlags(target *string, flags *pflag.FlagSet, flagName string) error {
	if flag := flags.Lookup(flagName); flag == nil {
		return nil
	}
	value, err := flags.GetString(flagName)
	if err != nil {
		return err
	}
	if value == "" {
		return nil
	}
	if *target != "" {
		return duplicateConfigErrors[flagName]
	}
	*target = value
	return nil
}

// DocumentFromFlags sets t.DocID from the passed flagset.
func (t *Target) DocumentFromFlags(flags *pflag.FlagSet) error {
	return setFromFlags(&t.Document, flags, FlagDocument)
}

// DatabaseFromFlags sets t.Database from the passed flagset.
func (t *Target) DatabaseFromFlags(flags *pflag.FlagSet) error {
	return setFromFlags(&t.Database, flags, FlagDatabase)
}

// FilenameFromFlags sets t.Filename from the passed flagset.
func (t *Target) FilenameFromFlags(flags *pflag.FlagSet) error {
	return setFromFlags(&t.Filename, flags, FlagFilename)
}

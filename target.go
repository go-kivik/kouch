package kouch

import (
	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch/internal/errors"
	"github.com/spf13/pflag"
)

// TargetScope represents the scope for a target, as relative targets have different
// meanings in different contexts.
type TargetScope int

// The supported target scopes
const (
	TargetRoot TargetScope = iota
	TargetDatabase
	TargetDocument
	TargetAttachment
	// View
	// Show
	// List
	// Update
	// Rewrite ??
	TargetLastScope = iota - 1
)

var _ = targetLastScope // lastScope only use in tests; this prevents linter warnings

// TargetScopeName returns the name of the scope, or "" if scope is invalid.
func TargetScopeName(scope TargetScope) string {
	switch scope {
	case TargetRoot:
		return "root"
	case TargetDatabase:
		return "database"
	case TargetDocument:
		return "document"
	case TargetAttachment:
		return "attachment"
	}
	return ""
}

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

// NewClient returns a chttp.Client, connected to the target server
func (t *Target) NewClient() (*chttp.Client, error) {
	if t.Root == "" {
		return nil, errors.NewExitError(chttp.ExitFailedToInitialize, "no server root specified")
	}
	c, err := chttp.New(t.Root)
	if err != nil {
		return nil, err
	}
	c.UserAgents = append(c.UserAgents, "Kouch/"+Version)
	if t.Username != "" || t.Password != "" {
		return c, c.Auth(&chttp.BasicAuth{
			Username: t.Username,
			Password: t.Password,
		})
	}
	return c, nil
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

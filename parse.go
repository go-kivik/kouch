package kouch

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch/internal/errors"
	"github.com/spf13/pflag"
)

// Target is a parsed target passed on the command line
type Target struct {
	// Root is the root URL.
	Root string
	// Database is the database name.
	Database string
	// DocID is the document ID.
	DocID string
	// Filename is the attachment filename.
	Filename string
}

// ParseAttachmentTarget parses a target containing a possible attachment ID.
func ParseAttachmentTarget(target string) (*Target, error) {
	if strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") {
		url, err := url.Parse(target)
		if err != nil {
			return nil, &errors.ExitError{Err: err, ExitCode: chttp.ExitStatusURLMalformed}
		}
		tgt, err := ParseAttachmentTarget(url.Path)
		tgt.Root = fmt.Sprintf("%s://%s/", url.Scheme, url.Host)
		return tgt, err
	}
	if strings.HasPrefix(target, "/") {
		parts := strings.SplitN(target, "/", 4)
		if len(parts) < 4 {
			return nil, errors.NewExitError(chttp.ExitStatusURLMalformed, "invalid target")
		}
		return &Target{
			Database: parts[1],
			DocID:    parts[2],
			Filename: parts[3],
		}, nil
	}
	if strings.Contains(target, "/") {
		parts := strings.SplitN(target, "/", 2)
		return &Target{
			DocID:    parts[0],
			Filename: parts[1],
		}, nil
	}
	return &Target{Filename: target}, nil
}

// FilenameFromFlags sets t.Filename from the passed flagset.
func (t *Target) FilenameFromFlags(flags *pflag.FlagSet) error {
	fn, err := flags.GetString(FlagFilename)
	if err != nil {
		return err
	}
	if fn == "" {
		return nil
	}
	if t.Filename != "" {
		return &errors.ExitError{
			Err:      errors.New("Must not use --" + FlagFilename + " and pass separate filename"),
			ExitCode: chttp.ExitFailedToInitialize,
		}
	}
	t.Filename = fn
	return nil
}

// DocIDFromFlags sets t.DocID from the passed flagset.
func (t *Target) DocIDFromFlags(flags *pflag.FlagSet) error {
	id, err := flags.GetString(FlagDocID)
	if err != nil {
		return err
	}
	if id == "" {
		return nil
	}
	if t.DocID != "" {
		return &errors.ExitError{
			Err:      errors.New("Must not use --" + FlagDocID + " and pass document ID as part of the target"),
			ExitCode: chttp.ExitFailedToInitialize,
		}
	}
	t.DocID = id
	return nil
}

// DatabaseFromFlags sets t.Database from the passed flagset.
func (t *Target) DatabaseFromFlags(flags *pflag.FlagSet) error {
	db, err := flags.GetString(FlagDatabase)
	if err != nil {
		return err
	}
	if db == "" {
		return nil
	}
	if t.Database != "" {
		return &errors.ExitError{
			Err:      errors.New("Must not use --" + FlagDatabase + " and pass database as part of the target"),
			ExitCode: chttp.ExitFailedToInitialize,
		}
	}
	t.Database = db
	return nil
}

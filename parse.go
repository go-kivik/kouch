package kouch

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch/internal/errors"
	"github.com/spf13/pflag"
)

const (
	FlagFilename = "filename"
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

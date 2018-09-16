package kouch

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch/internal/errors"
	"github.com/spf13/pflag"
)

var errIncompleteURL = errors.NewExitError(chttp.ExitFailedToInitialize, "incomplete target URL")

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
	targetLastScope = iota - 1
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

// NewTarget builds a new target from the context and flags.
func NewTarget(ctx context.Context, scope TargetScope, flags *pflag.FlagSet) (*Target, error) {
	t := &Target{}

	if tgt := GetTarget(ctx); tgt != "" {
		var err error
		t, err = ParseTarget(scope, tgt)
		if err != nil {
			return nil, err
		}
	}

	if err := t.FilenameFromFlags(flags); err != nil {
		return nil, err
	}
	if err := t.DocumentFromFlags(flags); err != nil {
		return nil, err
	}
	if err := t.DatabaseFromFlags(flags); err != nil {
		return nil, err
	}

	if defCtx, err := Conf(ctx).DefaultCtx(); err == nil {
		if t.Root == "" {
			t.Root = defCtx.Root
		}
	}

	return t, nil
}

// ParseTarget parses src as a CouchDB target, according to the rules for scope.
func ParseTarget(scope TargetScope, src string) (*Target, error) {
	target := &Target{}
	if src == "" {
		return target, nil
	}
	if strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://") {
		url, err := url.Parse(src)
		if err != nil {
			return nil, errors.WrapExitError(chttp.ExitStatusURLMalformed, err)
		}
		src = url.EscapedPath()
		target.Root = fmt.Sprintf("%s://%s", url.Scheme, url.Host)
		target.Username = url.User.Username()
		target.Password, _ = url.User.Password()
	}
	switch scope {
	case TargetRoot:
		return root(target, src)
	case TargetDatabase:
		return database(target, src)
	case TargetDocument:
		return document(target, src)
	case TargetAttachment:
		return attachment(target, src)
	}
	return nil, errors.New("invalid scope")
}

func root(t *Target, src string) (*Target, error) {
	t.Root = t.Root + src
	return t, validate(t)
}

func database(t *Target, src string) (*Target, error) {
	src, t.Database = lastSegment(src)
	if t.Database == "" && t.Document == "" && t.Filename == "" {
		return nil, errIncompleteURL
	}
	return root(t, src)
}

func document(t *Target, src string) (*Target, error) {
	src, t.Document = chopDocument(src)
	if t.Document == "" && t.Filename == "" {
		return nil, errIncompleteURL
	}
	return database(t, src)
}

func attachment(t *Target, src string) (*Target, error) {
	src, t.Filename = lastSegment(src)
	if t.Filename == "" {
		return nil, errIncompleteURL
	}
	return document(t, src)
}

func lastSegment(src string) (string, string) {
	parts := strings.Split(src, "/")
	l := len(parts)
	return strings.Join(parts[0:l-1], "/"), parts[l-1]
}

func validate(t *Target) error {
	parts := []string{t.Root, t.Database, t.Document, t.Filename}
	test := strings.Trim(strings.Join(parts, "\t"), "\t")
	if strings.Contains(test, "\t\t") {
		// This means one of the inner elements is empty
		return errIncompleteURL
	}
	return nil
}

// chopDocument chops the document ID off the right end of the string, returning
// the two segments.
func chopDocument(src string) (string, string) {
	parts := strings.Split(src, "/")
	l := len(parts)
	if l > 1 && (parts[l-2] == "_design" || parts[l-2] == "_local") {
		return strings.Join(parts[0:l-2], "/"), strings.Join(parts[l-2:], "/")
	}
	return strings.Join(parts[0:l-1], "/"), parts[l-1]
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

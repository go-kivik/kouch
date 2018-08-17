package target

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch/internal/errors"
)

// Scope represents the scope for a target, as relative targets have different
// meanings in different contexts.
type Scope int

// The supported target scopes
const (
	Root Scope = iota
	Database
	Document
	Attachment
	// View
	// Show
	// List
	// Update
	// Rewrite ??
	lastScope = iota - 1
)

// ScopeName returns the name of the scope, or "" if scope is invalid.
func ScopeName(scope Scope) string {
	switch scope {
	case Root:
		return "root"
	case Database:
		return "database"
	case Document:
		return "document"
	case Attachment:
		return "attachment"
	}
	return ""
}

// Target is a parsed target passed on the command line
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

var errIncompleteURL = errors.NewExitError(chttp.ExitFailedToInitialize, "incomplete target URL")

func (t *Target) validate() error {
	parts := []string{t.Root, t.Database, t.Document, t.Filename}
	test := strings.Trim(strings.Join(parts, "\t"), "\t")
	if strings.Contains(test, "\t\t") {
		// This means one of the inner elements is empty
		return errIncompleteURL
	}
	return nil
}

func (t *Target) root(src string) (*Target, error) {
	t.Root = t.Root + src
	return t, t.validate()
}

func (t *Target) database(src string) (*Target, error) {
	src, t.Database = lastSegment(src)
	if t.Database == "" && t.Document == "" && t.Filename == "" {
		return nil, errIncompleteURL
	}
	return t.root(src)
}

func (t *Target) document(src string) (*Target, error) {
	src, t.Document = chopDocument(src)
	if t.Document == "" && t.Filename == "" {
		return nil, errIncompleteURL
	}
	return t.database(src)
}

func (t *Target) attachment(src string) (*Target, error) {
	src, t.Filename = lastSegment(src)
	if t.Filename == "" {
		return nil, errIncompleteURL
	}
	return t.document(src)
}

// Parse parses src as a CouchDB target, according to the rules for scope.
func Parse(scope Scope, src string) (*Target, error) {
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
	case Root:
		return target.root(src)
	case Database:
		return target.database(src)
	case Document:
		return target.document(src)
	case Attachment:
		return target.attachment(src)
	}
	return nil, errors.New("invalid scope")
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

func lastSegment(src string) (string, string) {
	parts := strings.Split(src, "/")
	l := len(parts)
	return strings.Join(parts[0:l-1], "/"), parts[l-1]
}

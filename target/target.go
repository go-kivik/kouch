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
}

// Parse parses src as a CouchDB target, according to the rules for scope.
func Parse(scope Scope, src string) (*Target, error) {
	if src == "" {
		return &Target{}, nil
	}
	switch scope {
	case Root:
		return &Target{Root: src}, nil
	case Database:
		if strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://") {
			url, err := url.Parse(src)
			if err != nil {
				return nil, errors.WrapExitError(chttp.ExitStatusURLMalformed, err)
			}
			root, db := lastSegment(url.Path)
			return &Target{
				Root:     fmt.Sprintf("%s://%s%s", url.Scheme, url.Host, root),
				Database: db,
			}, nil
		}
		root, db := lastSegment(src)
		return &Target{
			Root:     root,
			Database: db,
		}, nil
	case Document:
		if strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://") {
			url, err := url.Parse(src)
			if err != nil {
				return nil, errors.WrapExitError(chttp.ExitStatusURLMalformed, err)
			}
			root, doc := chopDocument(url.Path)
			root, db := lastSegment(root)
			if doc == "" || db == "" {
				return nil, errors.NewExitError(chttp.ExitFailedToInitialize, "incomplete target URL")
			}
			return &Target{
				Root:     fmt.Sprintf("%s://%s%s", url.Scheme, url.Host, root),
				Database: db,
				Document: doc,
			}, nil
		}
		db, doc := chopDocument(src)
		root, db := lastSegment(db)
		return &Target{
			Root:     root,
			Database: db,
			Document: doc,
		}, nil
	case Attachment:
		if strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://") {
			url, err := url.Parse(src)
			if err != nil {
				return nil, errors.WrapExitError(chttp.ExitStatusURLMalformed, err)
			}
			doc, att := lastSegment(url.Path)
			root, doc := chopDocument(doc)
			root, db := lastSegment(root)
			if att == "" || doc == "" || db == "" {
				return nil, errors.NewExitError(chttp.ExitFailedToInitialize, "incomplete target URL")
			}
			return &Target{
				Root:     fmt.Sprintf("%s://%s%s", url.Scheme, url.Host, root),
				Database: db,
				Document: doc,
				Filename: att,
			}, nil
		}
		doc, att := lastSegment(src)
		db, doc := chopDocument(doc)
		root, db := lastSegment(db)
		return &Target{
			Root:     root,
			Database: db,
			Document: doc,
			Filename: att,
		}, nil
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
	return strings.Join(parts[0:len(parts)-1], "/"),
		parts[len(parts)-1]
}

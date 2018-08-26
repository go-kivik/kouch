package util

import (
	"fmt"
	"net/url"

	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch"
)

// DocPath calculates the server path to a document
func DocPath(o *kouch.Options) string {
	return fmt.Sprintf("/%s/%s", url.QueryEscape(o.Database), chttp.EncodeDocID(o.Document))
}

// AttPath calculates the server path to an attachment.
func AttPath(o *kouch.Options) string {
	return fmt.Sprintf("/%s/%s/%s", url.QueryEscape(o.Database), chttp.EncodeDocID(o.Document), url.QueryEscape(o.Filename))
}

// DatabasePath calculates the server path to a database.
func DatabasePath(o *kouch.Options) string {
	return fmt.Sprintf("/%s", url.QueryEscape(o.Database))
}

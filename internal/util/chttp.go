package util

import (
	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch"
)

// NewChttp wraps chttp.New
func NewChttp(addr string) (*chttp.Client, error) {
	c, err := chttp.New(addr)
	if err != nil {
		return nil, err
	}
	c.UserAgents = append(c.UserAgents, "Kouch/"+kouch.Version)
	return c, nil
}

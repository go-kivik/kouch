package kouch

import "github.com/go-kivik/couchdb/chttp"

// Options represents the accumulated options passed by the user, through config
// files, the commandline, etc.
type Options struct {
	*Target
	*chttp.Options
}

// NewOptions returns a new, empty Options struct.
func NewOptions() *Options {
	return &Options{
		Target:  &Target{},
		Options: &chttp.Options{},
	}
}

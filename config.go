package kouch

import (
	"encoding/json"
	"fmt"
	"io"
)

// Config represents the kouch tool configuration.
type Config struct {
	// DefaultContext is the name of the context to be used by default.
	DefaultContext string `yaml:"default-context" json:"default-context,omitempty"`
	// Contexts is a map of referencable names to context configs
	Contexts []NamedContext `json:"contexts,omitempty"`

	// File is the file where config was read from, or more precisely, where
	// changes will be saved to.
	File string `json:"-" yaml:"-"`
}

// NamedContext relates nicknames to context information.
type NamedContext struct {
	// Name is the nickname for this Context
	Name string `json:"name"`
	// Context holds the context information
	Context *Context `json:"context"`
}

// Context is a server context (URL, auth info, session store, etc)
type Context struct {
	// Root is the URL to the server's root.
	Root     string `json:"root"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
}

// DefaultCtx returns the default context.
func (c *Config) DefaultCtx() (*Context, error) {
	name := c.DefaultContext
	if name == "" {
		return nil, InitError("No default context")
	}
	for _, nc := range c.Contexts {
		if nc.Name == name {
			return nc.Context, nil
		}
	}
	return nil, InitError(fmt.Sprintf("Default context '%s' not defined", name))
}

// Dump dumps the config as a JSON string on r. Any errors will be returned as
// an error on r.Read().
func (c *Config) Dump() (r io.ReadCloser) {
	r, w := io.Pipe()
	go func() {
		err := json.NewEncoder(w).Encode(c)
		_ = w.CloseWithError(err)
	}()
	return r
}

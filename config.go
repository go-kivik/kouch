package kouch

import (
	"fmt"
)

// Config represents the kouch tool configuration.
type Config struct {
	// DefaultContext is the name of the context to be used by default.
	DefaultContext string `yaml:"default-context"`
	// Contexts is a map of referencable names to context configs
	Contexts []NamedContext
}

// NamedContext relates nicknames to context information.
type NamedContext struct {
	// Name is the nickname for this Context
	Name string
	// Context holds the context information
	Context *Context
}

// Context is a server context (URL, auth info, session store, etc)
type Context struct {
	// Root is the URL to the server's root.
	Root string
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

package kouch

// Config represents the kouch tool configuration.
type Config struct {
	// DefaultContext is the name of the context to be used by default.
	DefaultContext string
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

package kouch

import (
	"context"
	"io"
	"sync"

	"github.com/spf13/cobra"
)

type contextKey struct {
	name string
}

// Context Keys
var (
	verboseContextKey         = &contextKey{"verbose"}
	outputContextKey          = &contextKey{"output"}
	outputProcessorContextKey = &contextKey{"outputer"}
	configContextKey          = &contextKey{"config"}
)

// CmdContext is the command execution context.
type CmdContext struct{}

// Conf returns the context's current configuration struct, or panics if none is
// set.
func Conf(ctx context.Context) *Config {
	return ctx.Value(configContextKey).(*Config)
}

// SetConf returns a new context with the current config set to conf.
func SetConf(ctx context.Context, conf *Config) context.Context {
	return context.WithValue(ctx, configContextKey, conf)
}

// Outputer returns the context's current output processor, or panics if none
// is set.
func Outputer(ctx context.Context) OutputProcessor {
	return ctx.Value(outputProcessorContextKey).(OutputProcessor)
}

// SetOutputer returns a new context with the output processor set to op.
func SetOutputer(ctx context.Context, op OutputProcessor) context.Context {
	return context.WithValue(ctx, outputProcessorContextKey, op)
}

// Output returns the context's current output, or panics if none is set.
func Output(ctx context.Context) io.Writer {
	return ctx.Value(outputContextKey).(io.Writer)
}

// SetOutput returns a new context with the output set to w.
func SetOutput(ctx context.Context, w io.Writer) context.Context {
	return context.WithValue(ctx, outputContextKey, w)
}

// Verbose returns the verbosity flag of the context.
func Verbose(ctx context.Context) bool {
	verbose, _ := ctx.Value(verboseContextKey).(bool)
	return verbose
}

// SetVerbose returns a new context with the Verbose flag set to value.
func SetVerbose(ctx context.Context, value bool) context.Context {
	return context.WithValue(ctx, verboseContextKey, value)
}

type contexter interface {
	Context() context.Context
}

// GetContext returns the context associated with cmd.
func GetContext(cmd *cobra.Command) context.Context {
	// First, check if my PR (https://github.com/spf13/cobra/pull/727) has
	// been merged...
	if cxer, ok := interface{}(cmd).(contexter); ok {
		return cxer.Context()
	}
	return getContext(cmd)
}

var contexts map[*cobra.Command]context.Context
var contextMU = new(sync.RWMutex)

// getContext uses an ugly hack, inspired by Gorilla's contexts, to associate
// a context with a specific *cobra.Context instance. These instances are never
// cleaned up, but for a CLI app, that seems okay--typically during normal
// usage, there will be only one. And in tests, the processes are also
// short-lived, so waiting for the process to exit shouldn't be an issue in
// practice.
//
// If there is no context associated with the command, one is created from
// context.Background(), assgigned to the map, and returned.
//
// If my PR (https://github.com/spf13/cobra/pull/727), or an equivalent, is
// ever merged, this hack can be eliminated.
func getContext(cmd *cobra.Command) context.Context {
	initContext()
	contextMU.RLock()
	ctx, ok := contexts[cmd]
	contextMU.RUnlock()
	if ok {
		return ctx
	}
	ctx = context.Background()
	setContext(ctx, cmd)
	return ctx
}

type contextSetter interface {
	SetContext(context.Context)
}

// SetContext sets the context associated with the command.
func SetContext(ctx context.Context, cmd *cobra.Command) {
	// First, check if my PR (https://github.com/spf13/cobra/pull/727) has
	// been merged...
	if cxer, ok := interface{}(cmd).(contextSetter); ok {
		cxer.SetContext(ctx)
	}
	setContext(ctx, cmd)
}

func setContext(ctx context.Context, cmd *cobra.Command) {
	initContext()
	contextMU.Lock()
	defer contextMU.Unlock()
	contexts[cmd] = ctx
}

func initContext() {
	if contexts != nil {
		return
	}
	contextMU.Lock()
	defer contextMU.Unlock()
	contexts = make(map[*cobra.Command]context.Context, 1)
}

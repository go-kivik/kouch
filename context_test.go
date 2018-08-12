package kouch

import (
	"context"
	"testing"

	"github.com/spf13/cobra"
)

func TestGetContext(t *testing.T) {
	t.Run("first run", func(t *testing.T) {
		cmd := &cobra.Command{Use: "a"}
		ctx := getContext(cmd)
		if ctx == nil {
			t.Fatal("getContext returned nil")
		}
		if len(contexts) != 1 {
			t.Fatalf("Expected exactly 1 context in map, found %d", len(contexts))
		}
	})
	t.Run("re-fetch", func(t *testing.T) {
		cmd1 := &cobra.Command{Use: "b"}
		cmd2 := &cobra.Command{Use: "c"}
		ctx1 := getContext(cmd1)
		key := contextKey{"foo"}
		ctx1 = context.WithValue(ctx1, key, int(123))
		setContext(ctx1, cmd1)
		ctx2 := getContext(cmd2)
		ctx1 = getContext(cmd1)
		if v, _ := ctx1.Value(key).(int); v != 123 {
			t.Error("Expected 123 value from ctx1")
		}
		if _, ok := ctx2.Value(key).(int); ok {
			t.Error("Expected not to find value in ctx2")
		}
	})
}

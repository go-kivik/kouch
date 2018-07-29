package testy

import (
	"os"
	"testing"

	"github.com/flimzy/diff"
)

func TestRestoreEnv(t *testing.T) {
	os.Clearenv()
	os.Setenv("foo", "bar")
	os.Setenv("bar", "baz")

	func() {
		defer RestoreEnv()()
		os.Setenv("baz", "qux")
		os.Unsetenv("bar")
		env := os.Environ()
		expected := []string{"foo=bar", "baz=qux"}
		if d := diff.Interface(expected, env); d != nil {
			t.Fatal(d)
		}
	}()

	env := os.Environ()
	expected := []string{"foo=bar", "bar=baz"}
	if d := diff.Interface(expected, env); d != nil {
		t.Fatal(d)
	}
}

func TestEnviron(t *testing.T) {
	os.Clearenv()
	os.Setenv("foo", "bar")
	os.Setenv("bar", "baz")
	expected := map[string]string{
		"foo": "bar",
		"bar": "baz",
	}
	env := Environ()
	if d := diff.Interface(expected, env); d != nil {
		t.Fatal(d)
	}
}

func TestSetEnv(t *testing.T) {
	os.Clearenv()
	os.Setenv("foo", "bar")

	SetEnv(map[string]string{
		"bar": "baz",
		"baz": "qux",
	})

	expected := map[string]string{
		"foo": "bar",
		"bar": "baz",
		"baz": "qux",
	}
	env := Environ()
	if d := diff.Interface(expected, env); d != nil {
		t.Fatal(d)
	}

}

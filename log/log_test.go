package log

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/flimzy/testy"
)

func TestOpenLogFile(t *testing.T) {
	type Test struct {
		name   string
		file   string
		force  bool
		err    string
		finish func(*testing.T, *os.File)
	}
	tests := []Test{
		{
			name: "no filename",
			file: "",
			finish: func(t *testing.T, f *os.File) {
				if f != nil {
					t.Fatal("returned file should be nil")
				}
			},
		},
		func() Test {
			dir, err := ioutil.TempDir("", "kouch")
			if err != nil {
				t.Fatal(err)
			}
			return Test{
				name: "file does not exist",
				file: dir + "/doesntexist",
				finish: func(t *testing.T, f *os.File) {
					if err := f.Close(); err != nil {
						t.Error(err)
					}
					if err := os.RemoveAll(dir); err != nil {
						t.Fatal(err)
					}
				},
			}
		}(),
		func() Test {
			file, err := ioutil.TempFile("", "kouch")
			if err != nil {
				t.Fatal(err)
			}
			_ = file.Close()
			return Test{
				name: "file does exist",
				file: file.Name(),
				finish: func(t *testing.T, f *os.File) {
					if f != nil {
						t.Errorf("Expected nil file")
					}
					if err := os.Remove(f.Name()); err != nil {
						t.Fatal(err)
					}
				},
				err: fmt.Sprintf("open %s: file exists", file.Name()),
			}
		}(),
		func() Test {
			file, err := ioutil.TempFile("", "kouch")
			if err != nil {
				t.Fatal(err)
			}
			_ = file.Close()
			return Test{
				name:  "force file does exist",
				file:  file.Name(),
				force: true,
				finish: func(t *testing.T, f *os.File) {
					if err := f.Close(); err != nil {
						t.Error(err)
					}
					if err := os.Remove(f.Name()); err != nil {
						t.Fatal(err)
					}
				},
			}
		}(),
	}
	for _, test := range tests {
		func(test Test) {
			t.Run(test.name, func(t *testing.T) {
				t.Parallel()
				f, err := OpenLogFile(test.file, test.force)
				testy.Error(t, test.err, err)
				if test.finish != nil {
					test.finish(t, f)
				}
			})
		}(test)
	}
}

package log

import (
	"fmt"
	"io"
	"os"
)

// Logger handles logging.
type Logger interface {
	SetVerbose(bool)
	Printf(format string, a ...interface{})
	Println(...interface{})
	Debugln(...interface{})
	Errorln(...interface{})
}

type log struct {
	verbose bool
	stdout  io.Writer
	stderr  io.Writer
}

var _ Logger = &log{}

// New returns a new Logger.
func New() Logger {
	return &log{
		stdout: os.Stdout,
		stderr: os.Stderr,
	}
}

func (l *log) SetVerbose(v bool) {
	l.verbose = v
}

func (l *log) SetStdout(out io.Writer) {
	l.stdout = out
}

func (l *log) SetStderr(out io.Writer) {
	l.stderr = out
}

// Println wraps fmt.Println
func (l *log) Debugln(a ...interface{}) {
	if !l.verbose {
		return
	}
	fmt.Fprintln(l.stderr, a...)
}

func (l *log) Println(a ...interface{}) {
	fmt.Fprintln(l.stdout, a...)
}

func (l *log) Printf(format string, a ...interface{}) {
	fmt.Fprintf(l.stdout, format, a...)
}

func (l *log) Errorln(a ...interface{}) {
	fmt.Fprintln(l.stderr, a...)
}

// OpenLogFile opens file for writing, truncating any existing file if force
// is true.
func OpenLogFile(file string, force bool) (*os.File, error) {
	if file == "" {
		return nil, nil
	}
	if force {
		return os.Create(file)
	}
	return os.OpenFile(file, os.O_CREATE|os.O_EXCL, 0666)
}

package log

import "fmt"

// Logger handles logging.
type Logger interface {
	SetVerbose(bool)
	Println(...interface{}) (int, error)
}

type log struct {
	verbose bool
}

var _ Logger = &log{}

// New returns a new Logger.
func New() Logger {
	return &log{}
}

func (l *log) SetVerbose(v bool) {
	l.verbose = v
}

// Println wraps fmt.Println
func (l *log) Println(a ...interface{}) (int, error) {
	if !l.verbose {
		return 0, nil
	}
	return fmt.Println(a...)
}

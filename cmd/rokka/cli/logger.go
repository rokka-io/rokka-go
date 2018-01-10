package cli

import (
	"fmt"
	"os"
)

// Log is a simple abstraction over stdout/stderr logging.
type Log struct {
	Verbose bool
}

// Errorf prints to os.Stderr
func (l Log) Errorf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
}

// Printf prints to os.Stdout
func (l Log) Printf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stdout, format, a...)
}

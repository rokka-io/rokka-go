package cli

import (
	"fmt"
	"os"
)

// cliLog is a simple abstraction over stdout/stderr logging.
type cliLog struct {
	Verbose bool
	StdErr  *os.File
	StdOut  *os.File
}

// newCLILog creates a cliLog with stdout/stderr set to default.
func newCLILog(verbose bool) *cliLog {
	return &cliLog{
		Verbose: verbose,
		StdErr:  os.Stderr,
		StdOut:  os.Stdout,
	}
}

// Errorf prints to os.Stderr
func (l cliLog) Errorf(format string, a ...interface{}) {
	fmt.Fprintf(l.StdErr, format, a...)
}
func (l cliLog) Error(str string) {
	fmt.Fprint(l.StdErr, str)
}

// Printf prints to os.Stdout
func (l cliLog) Printf(format string, a ...interface{}) {
	fmt.Fprintf(l.StdOut, format, a...)
}

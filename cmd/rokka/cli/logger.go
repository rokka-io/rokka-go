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

// NewCLILog creates a cliLog with stdout/stderr set to default.
func NewCLILog(verbose bool) *cliLog {
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

// Printf prints to os.Stdout
func (l cliLog) Printf(format string, a ...interface{}) {
	fmt.Fprintf(l.StdOut, format, a...)
}

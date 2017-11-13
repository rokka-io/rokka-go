package cli

import (
	"fmt"
	"os"
)

type Log struct {
	Verbose bool
}

func (l Log) Errorf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
}

func (l Log) Printf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stdout, format, a...)
}

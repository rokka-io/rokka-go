package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/rokka-io/rokka-go/cli"
	"github.com/rokka-io/rokka-go/rokka"
)

var apiKey string
var apiAddress string
var verbose bool

func init() {
	flag.StringVar(&apiKey, "apiKey", "", "Optional API key")
	flag.StringVar(&apiAddress, "apiAddress", "", "Optional API address")
	flag.BoolVar(&verbose, "verbose", false, "Verbose (output request/response)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <command>\n\n", os.Args[0])
		fmt.Fprint(os.Stderr, "Commands:\n")
		printCommands(os.Stderr)
		fmt.Fprint(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		flag.Usage()
		os.Exit(0)
	}

	cfg, err := cli.GetConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading configuration file: %s\n", err)
		os.Exit(1)
	}

	if len(apiKey) != 0 {
		cfg.APIKey = apiKey
	}

	if len(apiAddress) != 0 {
		cfg.APIAddress = apiAddress
	}

	cl := rokka.NewClient(&rokka.Config{
		APIKey:     cfg.APIKey,
		APIAddress: cfg.APIAddress,
		Verbose:    verbose,
	})

	err = cli.ExecCommand(cl, args)

	if err == nil {
		os.Exit(0)
	}

	switch err := err.(type) {
	case cli.UnknownCommandError:
		fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
		fmt.Fprint(os.Stderr, "Commands:\n")
		printCommands(os.Stderr)
	case rokka.StatusCodeError:
		fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
		fmt.Fprintf(os.Stderr, "%s\n", cli.PrettyJSON(err.Body))
	default:
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}

	os.Exit(1)
}

func printCommands(f *os.File) {
	for _, c := range cli.Commands {
		options := ""
		if len(c.Options) != 0 {
			options = fmt.Sprintf("\t (Options: %s)", strings.Join(c.Options, ", "))
		}
		fmt.Fprintf(f, "  %s\t%s%s\n", strings.Join(c.Args, " "), c.Description, options)
	}
}

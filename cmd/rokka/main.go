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
var logger cli.Log

func init() {
	logger = cli.Log{}

	flag.StringVar(&apiKey, "apiKey", "", "Optional API key")
	flag.StringVar(&apiAddress, "apiAddress", "", "Optional API address")
	flag.BoolVar(&verbose, "verbose", false, "Show verbose HTTP request/responses")

	flag.Usage = func() {
		logger.Errorf("Usage: %s <command>\n\n", os.Args[0])
		logger.Errorf("Commands:\n%s\n", getCommandUsages())
		logger.Errorf("%s", "\nOptions:\n")
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

	logger.Verbose = verbose

	hc := NewHTTPClient(&logger)

	cl := rokka.NewClient(&rokka.Config{
		APIKey:     cfg.APIKey,
		APIAddress: cfg.APIAddress,
		HTTPClient: hc,
	})

	err = cli.ExecCommand(cl, &logger, args)

	if err == nil {
		os.Exit(0)
	}

	switch err := err.(type) {
	case cli.UnknownCommandError:
		logger.Errorf("Error: %v\n\n", err)
		flag.Usage()
	case rokka.StatusCodeError:
		logger.Errorf("Error: %v\n\n", err)

		s, pErr := cli.PrettyJSON(err.Body)
		if pErr != nil {
			logger.Errorf("Error pretty printing JSON: %s", pErr)
		}
		logger.Printf("%s", s)
	default:
		logger.Errorf("Error: %v\n", err)
	}

	os.Exit(1)
}

func getCommandUsages() string {
	var s string
	for _, c := range cli.Commands {
		options := ""
		if len(c.Options) != 0 {
			options = fmt.Sprintf("\t (Options: %s)", strings.Join(c.Options, ", "))
		}
		s = s + fmt.Sprintf("  %s\t%s%s\n", strings.Join(c.Args, " "), c.Description, options)
	}

	return s
}

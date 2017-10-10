package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rokka-io/rokka-go/cli"
	"github.com/rokka-io/rokka-go/client"
)

var apiKey string

func init() {
	flag.StringVar(&apiKey, "apiKey", "", "Optional API Key")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <command>\n\n", os.Args[0])
		fmt.Fprint(os.Stderr, "Actions:\n")
		fmt.Fprint(os.Stderr, "Options:\n")
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
		fmt.Printf("Error reading configuration file: %s\n", err)
		os.Exit(1)
	}

	if len(apiKey) != 0 {
		cfg.APIKey = apiKey
	}

	cl := client.NewClient(&client.Config{
		APIKey: cfg.APIKey,
	})

	cli.ExecCommand(cl, args)
}

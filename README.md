# Rokka API client library for Go

[![GoDoc](https://godoc.org/github.com/rokka-io/rokka-go?status.svg)](https://godoc.org/github.com/rokka-io/rokka-go)
[![Build Status](https://travis-ci.org/rokka-io/rokka-go.svg?branch=master)](https://travis-ci.org/rokka-io/rokka-go)
[![Coverage Status](https://coveralls.io/repos/github/rokka-io/rokka-go/badge.svg?branch=master)](https://coveralls.io/github/rokka-io/rokka-go?branch=master)

*This is currently work in progress*

## CLI Installation & Usage

The rokka-go CLI is a single file binary, automatically built for Windows, macOS and Linux. [You can download the appropriate file on the releases page](https://github.com/rokka-io/rokka-go/releases/latest).
After downloading the binary, place it in a folder [within your $PATH](https://en.wikipedia.org/wiki/PATH_(variable)) and rename it to `rokka`. Ensure the executable flag is set (e.g. on macOS and Linux execute `chmod u+x rokka`).

After that you can start the CLI.

```
# show commands and general help
$ rokka

# show help about a command (e.g. `login`)
$ rokka help login

# execute a command (e.g. `stacks list`)
$ rokka stacks list <organization name>

# get the raw json value of a command (e.g. to pipe it to another command)
$ rokka stacks list <organization name> -r

# show verbose request and response header/body
$ rokka stacks list <organization name> -v
```

### Login

The `login` command can be used to store the API-Key. Without using `login` first, the API-Key has to be specified for all executed commands.

```
# configuration gets stored in $HOME/.rokka/config
$ rokka login --apiKey="ENTER-API-KEY-HERE"

# configuration gets stored in ./my-awesome-place
$ rokka login --apiKey="ENTER-API-KEY-HERE" --config=./my-awesome-place

# If you use a different location for storing your API key, don't forget to specify the config location for all executed commands
$ rokka stacks list <organization name> --config=./my-awesome-place
```

## Library Usage

The library can be imported using the package import path `github.com/rokka-io/rokka-go/rokka`.
The godoc is published on [godoc.org/github.com/rokka-io/rokka-go](https://godoc.org/github.com/rokka-io/rokka-go).

Basic usage example:

```
package main

import (
	"fmt"
	"os"

	"github.com/rokka-io/rokka-go/rokka"
)

func main() {
	c := rokka.NewClient(&rokka.Config{
		APIKey: "exampleAPIKey",
	})

	resp, err := c.GetOrganization("example")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("%s", resp.ID)
}
```

## Contributing

### Install

Go >=1.7 is required.

```sh
$ go get github.com/rokka-io/rokka-go
```

### Dependencies

This project uses [dep](https://github.com/golang/dep). Run `dep ensure` for dependencies.

### Development

```sh
# Run CLI during development
$ go run ./cmd/rokka/main.go

# Build platform specific executables
$ GOOS=darwin go build -o ./bin/rokka ./cmd/rokka
$ GOOS=linux go build -o ./bin/rokka ./cmd/rokka
$ GOOS=windows go build -o ./bin/rokka ./cmd/rokka

# Update (auto-generate) rokka/operations_objects.go
$ go generate ./rokka

# Run tests
$ go test ./...
```

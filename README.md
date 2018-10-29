# Rokka API client library for Go

[![GoDoc](https://godoc.org/github.com/rokka-io/rokka-go?status.svg)](https://godoc.org/github.com/rokka-io/rokka-go)
[![Build Status](https://travis-ci.org/rokka-io/rokka-go.svg?branch=master)](https://travis-ci.org/rokka-io/rokka-go)
[![Coverage Status](https://coveralls.io/repos/github/rokka-io/rokka-go/badge.svg?branch=master)](https://coveralls.io/github/rokka-io/rokka-go?branch=master)

## CLI Installation & Usage

The rokka-go CLI is a single file binary, automatically built for Windows, macOS and Linux. [You can download the appropriate file on the releases page](https://github.com/rokka-io/rokka-go/releases/latest).
After downloading the binary, place it in a folder [within your $PATH](https://en.wikipedia.org/wiki/PATH_(variable)) and rename it to `rokka`. Ensure the executable flag is set (e.g. on macOS and Linux execute `chmod u+x rokka`).

After that you can start the CLI.

```bash
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

```bash
# configuration gets stored in $HOME/.rokka/config
$ rokka login --apiKey="ENTER-API-KEY-HERE"

# configuration gets stored in ./my-awesome-place
$ rokka login --apiKey="ENTER-API-KEY-HERE" --config=./my-awesome-place

# If you use a different location for storing your API key, don't forget to specify the config location for all executed commands
$ rokka stacks list <organization name> --config=./my-awesome-place
```

### Advanced usage

The output of every command in the CLI is defined by a template written in Go's [text/template language](https://golang.org/pkg/text/template/).
In case you want to have a more specific output adapted to your needs you can overwrite that template on the fly by specifying the `--template` flag.

Example:
```bash
# list source images
$ rokka sourceimages list <organization>
Name                                           Hash                                      Details
rokka-pic.png                     AAAA937b9b057e419cf96c0696be8db9ed481BBB  image/png, 1260x840
foo.jpg                           AAAA3b1297cd6c272f5beb253921956b81007BBB  image/jpeg, 2498x1389

# list only the hashes (make sure you enter the newlines correctly)
$ rokka sourceimages list <organization> --template="
{{range .Items}}{{.Hash}}
{{end}}"

AAAA937b9b057e419cf96c0696be8db9ed481BBB
AAAA3b1297cd6c272f5beb253921956b81007BBB
```

## Library Usage

Go >=1.8 is required.

```bash
$ go get github.com/rokka-io/rokka-go
```

The library can be imported using the package import path `github.com/rokka-io/rokka-go/rokka`.
The godoc is published on [godoc.org/github.com/rokka-io/rokka-go](https://godoc.org/github.com/rokka-io/rokka-go).

Basic usage example:

```go
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

### Automatic retries

The library supports a mechanism to auto retry requests if they fail. The retries happen
if either a HTTP, network or transport error occurs, or a status code of
429 (too many requests), 502 (bad gateway), 503 (service unavailable), or 504 (gateway timeout) has been received.

You can enable auto retry only on certain requests or, by storing the return value of `AutoRetry()` for all requests.

```go
c := rokka.NewClient(&rokka.Config{
	APIKey: "exampleAPIKey",
})

// enable for one request
resp, err := c.AutoRetry().GetOrganization("example")

// store return value and use it for all requests
retryingClient := c.AutoRetry()
retryingClient.GetOrganization("example")
```

## Contributing

### Dependencies

This project uses [dep](https://github.com/golang/dep). Run `dep ensure` for dependencies.

### Development

```bash
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
